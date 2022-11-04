---
title: "Trino Overview"
date: 2022-11-04
draft: true
tags:
- trino
featured: false
authors: ["ngoctd"]
---
Trino is a distributed SQL query engine designed to query large data sets distributed over one or more heterogeneous data sources.

Since Trino is being called a database by many members of the community, it makes sense to begin with a definition of what Trino is not. Trino is not a general-purpose relation database. It is not a replacement for databases like MySQL, PostgreSQL or Oracle. 

## What Trino is

Trino is a tool designed to efficiently query vast amount of data using distributed queries. If you work with terabytes or petabytes of data, you are likely using tools that interact with Hadoop and HDFS.

Trino was designed as an alternative to tools that query HDFS using pipelines of MapReduce jobs, such as Hive or Pig, but Trino it not limited to HDFS. Trino can be and has been extended to operate over different kinds of data sources, including traditional relational databases and other data sources such as Cassandra.

Trino was designed to handle data warehousing and analytics: data analysis, aggregating large amount of data and producing reports. These workloads are often classified as OLAP.


## Trino component: coordinator ( server )

The Trino coordinator is the server that is responsible for parsing statements, planning queries, and managing Trino worker nodes. Every Trino installation must have at least one Trino coordinator alongside one or more Trino workers.
Coordinators communicate with workers and clients using REST API.

## Trino component: worker ( server )

The Trino worker is responsible for executing tasks and processing data. Worker nodes fetch data from connectors and exchange intermediate data with each other. The coordinator is responsible for fetching results from the workers and returning the final results to the client.

When a Trino worker process starts up, it advertises itself to the discovery server in the coordinator, which makes it available to the Trino coordinator for task execution.

Workers communicate with other workers and Trino coordinators using a REST API.

## Trino component: connector ( data source )

A connector adapts Trino to a data source such as Hive or a relational database. You can think of a connector the same way you think of a driver for a database. Connector implements Trino's SPI.

## Trino component: catalog ( data source )

## Trino component: schema ( data source )

## Hive Metadata Storage

TL;DR: The Hive connector is what you use in Trino for reading data from object storage that is organized according to the rules laid out by Hive, without using the Hive runtime code.

**References**

- [trino documentation](https://trino.io/docs/current/overview/use-cases.html)