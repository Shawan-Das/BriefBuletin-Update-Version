SHELL := pwsh.exe

.PHONY: fetch start stop status run

fetch:
	@git pull
	@py main.py
	@echo "Done!"
	@git add .
	@git commit -m "Fetch Done"
	@git push

# Start the background loop (runs every 30 minutes)
start:
	@git pull
	@pwsh -NoProfile -ExecutionPolicy Bypass -Command "& { .\scripts\run_loop.ps1 start }"

# Stop the background loop
stop:
	@pwsh -NoProfile -ExecutionPolicy Bypass -Command "& { .\scripts\run_loop.ps1 stop }"
	
# 	@git add .
# 	@git commit -m "Stopped Loop"
# 	@git push

# Show status (running or not)
status:
	@pwsh -NoProfile -ExecutionPolicy Bypass -Command "& { .\scripts\run_loop.ps1 status }"

# Run the worker loop in the foreground (for debugging)
run:
	@pwsh -NoProfile -ExecutionPolicy Bypass -File .\scripts\loop_worker.ps1