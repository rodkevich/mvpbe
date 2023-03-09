# create env file for settings
cp ./examples/example.env .env

docker-compose up -d --build --remove-orphans postgres
docker-compose up -d --build --remove-orphans rabbitmq
docker-compose up -d --build --remove-orphans migrator
