new_day:
	@read -p "Enter the day: " day; \
	cp templates/dayX solutions/day$$day -r; \
	mv solutions/day$$day/dayX.go solutions/day$$day/day$$day.go