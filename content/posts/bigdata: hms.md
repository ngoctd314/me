---
title: "Hive metastore service "
date: 2022-11-15
authors: ["ngoctd"]
draft: true
tags: ["bigdata", "trino"]
---

We can simplify the Hive architecture to four components:

- The runtime contains the logic of the query engine that translates the SQL - esque Hive Query Language (HQL) into MapReduce jobs that run over files stored in the filesystem.

- Storage component is simply that, it stores file in various formats and index structures to recall these files (JSON, CSV, ORC, Parquet, HDFS, Aws S3, GCS)

- In order for Hive to process these files, it must have a mapping from SQL tables in the runtime to files and directories in the storage component. To accomplish this, Hive uses the HMS (Hive Metastore Service)

## Hive connector

The hive connector allows querying data stored in an Apache Hive data warehouse. Hive is a combination of there components:

- Data files in varying formats (HDFS, S3)
- HMS
- Query language (HiveQL)

Trino only uses the first two components: the data and the metadata. It does not use HiveQL or any part of Hive's execution environment.