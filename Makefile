build-front:
	cd frontend && npm run build
	rm -rf isumid/handlers/front-built
	mv frontend/out isumid/handlers/front-built
	