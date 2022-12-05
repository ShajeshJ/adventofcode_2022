new_day:
	@read -p "Enter which day to create: " day; \
	cp templates/dayX solutions/day$$day -r; \
	mv solutions/day$$day/dayX.go solutions/day$$day/day$$day.go

run:
	@read -p "Enter which day to run: " day; \
	go run solutions/day$$day/day$$day.go