.PHONY: gen
gen:
	@mkdir -p logs
	@printf "%sDATABASE_NAME=second_db\nDATABASE_PORT=3307\nDD_API_KEY=\n" > .env
	@echo "\nSuccess!\nNow set DD_API_KEY in .env\n"
