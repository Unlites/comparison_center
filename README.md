# Comparison Center

Service for centralized and custom comparison of any categories of objects.

## Install

```shell
git clone https://github.com/Unlites/comparison_center
```

After cloning the project, create .env file based on .env.sample with overriding parameters if necessary.

Make sure that Docker service is currently running on your system. Then:

```shell
make run
make migrate_up
```

## Usage

Open your browser on http://localhost:3000 (or define another CLIENT_HOSTPORT at .env file).

You can create different comparisons with different custom options. After creating comparison, you can add object you're comparing, view objects you've already added, and sort them by rating, date added, and more.

## Metrics

You can visit http://localhost:3100 (or define another GRAFANA_HOSTPORT at .env file) and log into Grafana with admin:admin userpass. 

Unfortunately, Grafana does not support import dashboard with datasource already selected using Docker volumes since Grafana 9.x.x+, so you have to import it manually using UI.

To do this, import the dashboard from ./metrics/dashboards/requests.json with chosing Prometheus datasouce, which will be immediately able after installing. This dashboard provides you information about HTTP requests to application.

You can also add Node Exporter dashboard to monitor your system using the dashboard id 1860, as required service for that will be running too.

## Stack

 - Go
 - Vue.js
 - Docker
 - MongoDB
 - Migrations with [golang-migrate](https://github.com/golang-migrate/migrate)
 - Metrics with Prometheus and Grafana 

