.PHONY: gen
gen:
	@mkdir -p logs
	@touch logs/first.log logs/second.log
	@printf "%sDATABASE_NAME=second_db\nDATABASE_PORT=3307\nDD_API_KEY=\n" > .env
	@echo "\nSuccess!\nNow set DD_API_KEY in .env\n"
