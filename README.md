# Comparison Center

Service for centralized and custom comparison of any categories of objects.

## Install

```shell
git clone https://github.com/Unlites/comparison_center
make up
make migrate
```

## Usage

Open your browser on http://localhost:3000 (or define another CLIENT_PORT at .env file).

You can create different comparisons with different custom options. After creating comparison, you can add object you're comparing, view objects you've already added, and sort them by rating, date added, and more.

## Stack

 - Go
 - Vue.js
 - Docker
 - MongoDB
 - Migrations with [golang-migrate](https://github.com/golang-migrate/migrate)
 