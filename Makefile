goose_up:
	goose -dir migrations postgres "user=winte password=123 dbname=xployalty_test sslmode=disable" up

goose_down:
	goose -dir migrations postgres "user=winte password=123 dbname=xployalty_test sslmode=disable" down

goose_redo:
	goose -dir migrations postgres "user=winte password=123 dbname=xployalty_test sslmode=disable" redo

goose_status:
	goose -dir migrations postgres "user=winte password=123 dbname=xployalty_test sslmode=disable" status
