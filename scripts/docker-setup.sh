# create env file for settings
cp ./builders/docker/.env .env

docker-compose up --build --remove-orphans
