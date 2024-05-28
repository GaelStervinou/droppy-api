# DROPPY API

## INSTALLATION

- Create a .env file and fill it with th variables you can find in .env.example
- You can generate a random jwt secret by running `node -e "console.log(require('crypto').randomBytes(256).toString('base64'))"`
- Launch the command `docker compose build --no-cache && docker-compose up -d` or only `docker-compose up -d` if you have already built the image

You should have a functional API running on port 3000, with hot reload enabled.

## API DOCUMENTATION

You can find the API documentation at the following URL: [http://localhost:3000/api-docs](http://localhost:3000/swagger)
To generate the swagger documentation, run the following command: `swag init --parseDependency --parseInternal`