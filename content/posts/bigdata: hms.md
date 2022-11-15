---
title: "Hive meta storage "
date: 2022-11-15
authors: ["ngoctd"]
draft: true
tags: ["bigdata", "trino"]
---

We can simplify the Hive architecture to four components:

- The runtime contains the logic of the query engine that translates the SQL - esque Hive Query Language (HQL) into MapReduce jobs that run over files stored in the filesystem.

- Storage component is simply that, it stores file in various formats and index structures to recall these files (JSON, CSV, ORC, Parquet, HDFS, Aws S3, GCS)

- In order for Hive to process these files, it must have a mapping from SQL tables in the runtime to files and directories in the storage component. 
