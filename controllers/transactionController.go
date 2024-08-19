package controllers

import (
	"fmt"
	"math"
	"strconv"
	"sync"

	"go_api_tokoku/helpers"

	"github.com/gofiber/fiber/v2"
)

func GetCountTransactions(c *fiber.Ctx) error {

	customerId := c.Query("id")

	result, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT JSON_AGG(data.count_data) as count_data FROM (
														SELECT --ts.name as transaction_state, COUNT(tr.id) as count_data,
														JSONB_BUILD_OBJECT(ts.name, COUNT(tr.id)) as count_data
														FROM tk.transaction_state ts
														LEFT JOIN tk.transaction tr
															ON ts.id = tr.transaction_state_id
															AND tr.customer_id = %s
														GROUP BY ts.id
														ORDER BY ts.id
													) data`, customerId))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(result) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	flattenedResult := make(map[string]int)
	for _, item := range result {
		for _, value := range item {
			for _, val := range value.([]interface{}) {
				for k, v := range val.(map[string]interface{}) {
					flattenedResult[k] = int(v.(float64))
				}
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success",
		Success: true,
		Data:    flattenedResult,
	})

}

func executeGORMQuery(query string, resultsChan chan<- map[int][]map[string]interface{}, index int, wg *sync.WaitGroup) {
	defer wg.Done()

	results, _ := helpers.ExecuteQuery(query)

	resultsChan <- map[int][]map[string]interface{}{index: results}
}

func GetTransactions(c *fiber.Ctx) error {

	customerId := c.Query("id")
	types := c.Query("type")
	transactionId := c.Query("transactionId")
	table := c.Query("table")
	page := c.Query("page")
	pageSize := c.Query("pageSize")

	var qLimit, qPage string
	iPage, _ := strconv.Atoi(page)
	iPageSize, _ := strconv.Atoi(pageSize)

	if pageSize != "" {
		qLimit = " LIMIT " + pageSize
	} else {
		qLimit = " LIMIT 20"
		iPageSize = 20
	}

	if page == "" {
		iPage = 0
	} else {
		iPage = iPage - 1
	}

	tempQ := strconv.Itoa(iPage * iPageSize)
	qPage = " OFFSET " + tempQ

	var query string

	if transactionId == "" && table == "" {
		query = fmt.Sprintf(
			`SELECT sq.* FROM (
				SELECT tr.id, 
					JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'id', trd.id,
							'produk_id', p.id,
							'code', p.code,
							'name', p.name,
							'harga', trd.harga,
							'image', p.foto,
							'point', trd.point,
							'qty', trd.qty,
							'produk_satuan', JSONB_BUILD_OBJECT(
															'carton', ps.carton,
															'ball', ps.ball,
															'slof', ps.slof,
															'pack', ps.pack
														)
						) ORDER BY p.id
					) as products,
					CASE WHEN COUNT(trd.id) = 1 THEN NULL ELSE COUNT(trd.id) - 1 END as jumlah_produk,
					CASE WHEN rev.id IS NOT NULL THEN
						JSONB_BUILD_OBJECT(
									'rating', rev.rating,
									'description', rev.description,
									'photo', rev.photo
								)
							ELSE NULL END as review,
					tr.total_transaction as total_order,
					tr.transaction_date,
					ts.name as type,
					'transaction' as table
				FROM tk.transaction_state ts
				JOIN tk.transaction tr
					ON ts.id = tr.transaction_state_id
				JOIN tk.transaction_detail trd
					ON tr.id = trd.transaction_id
				JOIN produk p
					ON trd.produk_id = p.id
				LEFT JOIN produk_satuan ps
					ON p.satuan_id = ps.id
				LEFT JOIN tk.review rev
					ON tr.id = rev.order_id
					AND rev.order_item_id IS NULL
				WHERE tr.customer_id = %s AND ts.name = UPPER('%s') AND tr.penjualan_id IS NULL
				GROUP BY tr.id, rev.id, ts.id

				UNION 

				SELECT tr.id, 
					JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'id', pj.id,
							'produk_id', p.id,
							'code', p.code,
							'name', p.name,
							'harga', pd.harga,
							'image', p.foto,
							'point', 0,
							'qty', pd.jumlah,
							'produk_satuan', JSONB_BUILD_OBJECT(
															'carton', ps.carton,
															'ball', ps.ball,
															'slof', ps.slof,
															'pack', ps.pack
														)
						) ORDER BY p.id
					) as products,
					CASE WHEN COUNT(pd.id) = 1 THEN NULL ELSE COUNT(pd.id) - 1 END as jumlah_produk,
					CASE WHEN rev.id IS NOT NULL THEN
						JSONB_BUILD_OBJECT(
									'rating', rev.rating,
									'description', rev.description,
									'photo', rev.photo
								)
							ELSE NULL END as review,
					pj.total_penjualan as total_order,
					pj.tanggal_penjualan as transaction_date,
					ts.name as type,
					'penjualan' as table
				FROM tk.transaction_state ts
				JOIN tk.transaction tr
					ON ts.id = tr.transaction_state_id
					AND tr.penjualan_id IS NOT NULL
				JOIN penjualan pj
					ON tr.penjualan_id = pj.id
				JOIN penjualan_detail pd
					ON pj.id = pd.penjualan_id
				JOIN produk p
					ON pd.produk_id = p.id
				LEFT JOIN produk_satuan ps
					ON p.satuan_id = ps.id
				LEFT JOIN tk.review rev
					ON tr.id = rev.order_id
					AND rev.order_item_id IS NULL
				WHERE tr.customer_id = %s AND ts.name = UPPER('%s') AND tr.penjualan_id IS NOT NULL
				GROUP BY tr.id, pj.id, rev.id, ts.id
				) sq
				ORDER BY sq.transaction_date DESC`, customerId, types, customerId, types)
	} else {

		if transactionId != "" && table != "" {
			if table == "transaction" {

				query = fmt.Sprintf(
					`WITH data_details as (

						WITH data_review as (
							SELECT sq.order_id, MAX(sq.salesman_id) as salesman_id,
									JSONB_AGG(sq.value_salesman) FILTER (WHERE sq.value_salesman IS NOT NULL)->0 as value_salesman,
									JSONB_AGG(sq.value_transaction) FILTER (WHERE sq.value_transaction IS NOT NULL)->0 as value_transaction
							FROM (
							SELECT order_id, salesman_id,
									CASE WHEN order_id IS NOT NULL 
											AND salesman_id IS NOT NULL 
											THEN JSONB_BUILD_OBJECT(
												'id', rev.id,
												'rating', rev.rating,
												'description', rev.description,
												'photo', rev.photo
											)
																
										ELSE null
										END as value_salesman,
									CASE WHEN order_id IS NOT NULL
											AND salesman_id IS NULL
											THEN JSONB_BUILD_OBJECT(
												'id', rev.id,
												'rating', rev.rating,
												'description', rev.description,
												'photo', rev.photo
											)
																
										ELSE null
										END as value_transaction
									FROM tk.review rev WHERE order_id = %s
							) sq
							GROUP BY order_id
						)

					SELECT 
						trd.transaction_id,
						JSONB_AGG(
							JSONB_BUILD_OBJECT(
								'id', trd.id,
								'produk_id', trd.produk_id,
								'point', trd.point,
								'code', p.code,
								'name', p.name,
								'image', p.foto,
								'qty', trd.qty,
								'unit', iu.name,
								'harga', trd.harga,
								'note', trd.note,
								'produk_satuan', JSONB_BUILD_OBJECT(
													'carton', ps.carton,
													'ball', ps.ball,
													'slof', ps.slof,
													'pack', ps.pack
												)
							) ORDER BY p.order
						) as details,
						JSONB_BUILD_OBJECT(
							'id', tr.reference_id,
							'name', s.name,
							'phone', s.phone,
							'type', tr.reference_name
						) as courier,
						JSONB_BUILD_OBJECT(
							'id', c.id,
							'name', c.name,
							'outlet_name', c.outlet_name,
							'phone', c.phone,
							'address', c.alamat,
							'kelurahan', c.kelurahan,
							'kabupaten', c.kabupaten,
							'kecamatan', c.kecamatan,
							'provinsi', c.provinsi,
							'outlet_photo', c.outlet_photo
						) as customer,
						JSONB_BUILD_OBJECT(
							'no_transaksi', tr.id,
							'jumlah_pesanan', SUM(trd.qty),
							'jumlah_point', SUM(trd.point),
							'voucher_tokoku', 0,
							'total_pembayaran', tr.total_transaction
						) as invoice,
						COALESCE(rev.order_id, rev.salesman_id) as param_rev,
						JSONB_AGG(
							JSONB_BUILD_OBJECT(
								'salesman', rev.value_salesman,
								'transaction', rev.value_transaction
							)
						)->0 as review,
						comp.id as complaint_id,
						JSONB_BUILD_OBJECT(
							'description', comp.description,
							'photo', comp.image,
							'status', comp.status,
							'feedback', comp.feedback
						) as complaint
					FROM tk.transaction_detail trd
					JOIN tk.transaction tr
						ON trd.transaction_id = tr.id
					JOIN tk.transaction_state ts
						ON tr.transaction_state_id = ts.id
					JOIN customer c
						ON tr.customer_id = c.id
					JOIN produk p
						ON trd.produk_id = p.id
					JOIN produk_satuan ps
						ON p.satuan_id = ps.id
					LEFT JOIN tk.unit_mapping um
						ON p.id = um.produk_id
					LEFT JOIN tk.item_unit_tk iu
						ON um.item_unit_id = iu.id
					LEFT JOIN salesman s
						ON (tr.reference_id = s.id
						AND tr.reference_name = 'SALESMAN')
					LEFT JOIN data_review rev
						ON tr.id = rev.order_id
					LEFT JOIN tk.complaints comp
						ON tr.id = comp.transaction_id
						AND comp.transaction_item_id IS NULL
					WHERE tr.id = %s
					GROUP BY trd.transaction_id, tr.id, c.id, s.id, COALESCE(rev.order_id, rev.salesman_id), comp.id
				)
				
				SELECT LOWER(name) as name,
						JSONB_AGG(
							JSONB_BUILD_OBJECT(
								'id', tr.id,
								'type', ts.name,
								'transaction_date', tr.transaction_date,
								'total_order', tr.total_transaction,
								'products', trd.details,
								'courier', trd.courier,
								'customer', trd.customer,
								'invoice', trd.invoice,
								'note', tr.note,
								'review', CASE WHEN trd.param_rev IS NULL THEN NULL ELSE trd.review END,
								'complaints', CASE WHEN trd.complaint_id IS NULL THEN NULL ELSE trd.complaint END
							) ORDER BY tr.transaction_date DESC
						) FILTER (WHERE tr.id IS NOT NULL) as datas
				FROM tk.transaction_state ts
				LEFT JOIN tk.transaction tr
					ON ts.id = tr.transaction_state_id
				LEFT JOIN data_details trd
					ON tr.id = trd.transaction_id
				WHERE tr.id = %s
				GROUP BY ts.name`, transactionId, transactionId, transactionId)
			} else {

				query = fmt.Sprintf(`WITH data_penjualan as (
									SELECT 
										pd.penjualan_id,
										JSONB_AGG(
											JSONB_BUILD_OBJECT(
												'id', pj.id,
												'produk_id', pd.produk_id,
												'point', 0,
												'code', p.code,
												'name', p.name,
												'image', p.foto,
												'qty', pd.jumlah,
												'unit', 'Pack',
												'harga', pd.harga,
												'note', null,
												'produk_satuan', JSONB_BUILD_OBJECT(
																	'carton', ps.carton,
																	'ball', ps.ball,
																	'slof', ps.slof,
																	'pack', ps.pack
																)
											) ORDER BY p.order
										) as details,
										JSONB_BUILD_OBJECT(
											'id', s.id,
											'name', s.name,
											'phone', s.phone,
											'type', 'Salesman'
										) as courier,
										JSONB_BUILD_OBJECT(
											'id', c.id,
											'name', c.name,
											'outlet_name', c.outlet_name,
											'phone', c.phone,
											'address', c.alamat,
											'kelurahan', c.kelurahan,
											'kabupaten', c.kabupaten,
											'kecamatan', c.kecamatan,
											'provinsi', c.provinsi,
											'outlet_photo', c.outlet_photo
										) as customer,
										JSONB_BUILD_OBJECT(
											'no_transaksi', pj.id,
											'jumlah_pesanan', SUM(pd.jumlah),
											'jumlah_point', 0,
											'voucher_tokoku', 0,
											'total_pembayaran', pj.total_penjualan
										) as invoice,
										null as rating,
										null as review,
										null as complaint_id,
										null as complaint
									FROM penjualan_detail pd
									JOIN penjualan pj
										ON pd.penjualan_id = pj.id
									JOIN customer c
										ON pj.customer_id = c.id
									JOIN produk p
										ON pd.produk_id = p.id
									JOIN produk_satuan ps
										ON p.satuan_id = ps.id
									LEFT JOIN salesman s
										ON pj.salesman_id = s.id
									WHERE pj.id = %s
									GROUP BY pd.penjualan_id, pj.id, p.id, c.id, s.id
								)
									
								SELECT LOWER('Accepted') as name,
										JSONB_AGG(
											JSONB_BUILD_OBJECT(
												'id', tr.id,
												'type', 'Accepted',
												'transaction_date', tr.tanggal_penjualan,
												'total_order', tr.total_penjualan,
												'products', trd.details,
												'courier', trd.courier,
												'customer', trd.customer,
												'invoice', trd.invoice,
												'note', null,
												'review', CASE WHEN trd.rating IS NULL THEN NULL ELSE trd.review END,
												'complaints', CASE WHEN trd.complaint_id IS NULL THEN NULL ELSE trd.complaint END
											) ORDER BY tr.tanggal_penjualan DESC
										) FILTER (WHERE tr.id IS NOT NULL) as datas
								FROM penjualan tr
								LEFT JOIN data_penjualan trd
									ON tr.id = trd.penjualan_id
								WHERE tr.id = %s`, transactionId, transactionId)
			}
		}
	}

	results, err := helpers.ExecuteQuery(query)

	// go resultsCount, _ := helpers.ExecuteQuery(fmt.Sprintf(query + qPage + qLimit))
	// fmt.Println(query)

	var wg sync.WaitGroup
	resultsChan := make(chan map[int][]map[string]interface{}, 2)

	queries := []string{
		query,
		query + qPage + qLimit,
	}

	tempResults := make([][]map[string]interface{}, len(queries))

	// Launch concurrent Goroutines
	for i, query := range queries {
		wg.Add(1)
		go executeGORMQuery(query, resultsChan, i, &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		for index, res := range result {
			tempResults[index] = res
		}
	}

	if results == nil {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseDataMultiple{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	type newResponseDataMultiple struct {
		Message    string      `json:"message"`
		Success    bool        `json:"success"`
		Data       interface{} `json:"datas"`
		TotalPages int         `json:"total_pages"`
	}

	var tempTotalPages int
	if len(tempResults[0]) < iPageSize {
		tempTotalPages = 1
	} else {
		tempTotalPages = int(math.Ceil(float64(len(tempResults[0])) / float64(iPageSize)))
	}

	if transactionId != "" && table != "" {
		var newResults interface{}
		for _, val := range tempResults[1] {
			// fmt.Println(i, val)
			newResults = val["datas"]
		}
		return c.Status(fiber.StatusOK).JSON(newResponseDataMultiple{
			Message:    "Data has been loaded successfully",
			Success:    true,
			Data:       newResults,
			TotalPages: tempTotalPages,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(newResponseDataMultiple{
			Message:    "Data has been loaded successfully",
			Success:    true,
			Data:       tempResults[1],
			TotalPages: tempTotalPages,
		})

	}

}

func GetPointsCustomer(c *fiber.Ctx) error {

	customerId := c.Query("customerId")

	results, err := helpers.ExecuteQuery(fmt.Sprintf(
		`SELECT cph.customer_id,
				SUM(CASE WHEN cph.exchange_id IS NOT NULL THEN cph.point ELSE 0 END) as total_point_digunakan,
				SUM(CASE WHEN cph.transaction_id IS NOT NULL THEN cph.point ELSE 0 END) as total_point_terkumpul,
				(SUM(CASE WHEN cph.transaction_id IS NOT NULL THEN cph.point ELSE 0 END) - SUM(CASE WHEN cph.exchange_id IS NOT NULL THEN cph.point ELSE 0 END)) as total_point
		FROM tk.customer_point_history cph
		LEFT JOIN tk.transaction t
			ON cph.transaction_id = t.id
		WHERE TRUE AND cph.customer_id = %s
		GROUP BY cph.customer_id`, customerId))

	if results == nil {
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.Response{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    results[0],
	})
}

func GetPointsHistory(c *fiber.Ctx) error {

	customerId := c.Query("customerId")
	dateStart := c.Query("dateStart")
	dateEnd := c.Query("dateEnd")
	types := c.Query("type")

	whereQuery := ""
	if dateStart != "" && dateEnd != "" {
		whereQuery = " AND DATE(cph.datetime) BETWEEN DATE('" + dateStart + "') AND DATE('" + dateEnd + "')"
	}

	whereType := ""
	if types != "" {
		whereType = " AND UPPER(cph.type) = UPPER('" + types + "')"
	}

	results, err := helpers.ExecuteQuery(fmt.Sprintf(
		`SELECT cph.id,
				cph.customer_id,
				cph.transaction_id,
				cph.exchange_id,
				cph.datetime,
				COALESCE(cph.point,0) as point,
				cph.expired_date,
				cph.created_at,
				cph.updated_at, 
				t.transaction_date,
				cph.type,
				CASE 
					WHEN CURRENT_DATE > cph.expired_date THEN 'Expired at '||cph.expired_date
					WHEN cph.transaction_id IS NOT NULL THEN 'Transaction at '||t.transaction_date
				END as note
		FROM tk.customer_point_history cph
		LEFT JOIN tk.transaction t
		ON cph.transaction_id = t.id
		WHERE TRUE AND cph.customer_id = %s %s %s
		ORDER BY cph.datetime DESC`, customerId, whereQuery, whereType))

	if results == nil {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseDataMultiple{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    results,
	})
}
