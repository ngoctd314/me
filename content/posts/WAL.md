---
title: "Write Ahead Logging"
date: 2023-04-12T09:15:34+07:00
draft: true
authors: ["ngoctd"]
categories: ["concepts"]
--- 

Write-ahead logging (WAL) is a family of techniques for providing atomicity and durability (two of the ACID properties) in database systems. It can be seen as an implementation of the "Event Sourcing" architecture, in which the state of a system is the result of the evolution of incoming events from an initial state. A write ahead log is an append-only auxiliary disk-resident structure used for crash and transaction recovery. The changes are first recorded in the log, which must be written to stable storage, before the changes are written to the database.

The main functionality of a write-ahead log can be summarized as:
- Allow the page cache to buffer updates to disk-resident pages while ensuring durability semantics in the larger context of a database system.
- Persist all operations on disk until the cached copies of pages affect by these operations are synchronized on disk. Every operation that modifies state has to be logged on disk before the contents on the associated pages can be modified.
- Allow lost in-memory changes to be reconstructed from the operation log in case of a crash.

In a system using WAL, all modifications are written to a log before they are applied. Usually both redo and undo information is stored in the log.

A program is in the middle of performing some operation when the machine it is running on loses power. Upon restart, that program might need to know whether the operation it was performing succeeded, succeeded partially, or failed. If a write-ahead log is used, the program can check this log and compare what it was supposed to be doing when it unexpectedly lost power to what was actually done.

vòng techlead bị hỏi về oop, nginx, một số câu lệnh trên linux : check ram, cpu, hỏi về database: nosql(mongo), sql(postgre vs mysql), usecase sử dụng (nosql - sql), các cách tối ưu truy vấn (sử dụng index + partition) - ý nghĩa của 2 cái. Đưa ra 1 bài toán về thiết kế database (bài toán bán hàng). Ngoài ra hỏi về những dự án đã làm, luồng hoạt động của dự án

vòng pv với giám đốc khối (AI - platform): hỏi về threading, hệ điều hành, cách tăng tải server (scale instance + tách service theo nghiệp vụ), hỏi về cách hoạt động của kafka - so sánh nhẹ với rabitmq