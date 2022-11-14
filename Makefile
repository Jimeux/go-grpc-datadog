.PHONY: gen
gen:
	@mkdir -p logs
	@printf "%sDATABASE_NAME=second_db\nDATABASE_PORT=3307\nDD_API_KEY=\n" > .env
	@echo "\nSuccess!\nNow set DD_API_KEY in .env\n"

.PHONY: d-build-first d-build-second
d-build-second:
	docker build --build-arg SERVICE=second -t second-svc .
d-build-first:
	docker build --build-arg SERVICE=first -t first-svc .
