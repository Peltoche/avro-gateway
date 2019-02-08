[![license](http://img.shields.io/badge/license-Apache%20v2-orange.svg)](https://raw.githubusercontent.com/Peltoche/oaichecker/master/LICENSE)

# Avro Gateway

## What is it?

The Avro Gateway is a little API in front of a Schema Registry. It ensure that
all the consumers/producers schemas are compatible between them inside a topic.
It keeps track of all the schema versions used by the clients and refuse to serve
any schema versions susceptible to breaking the current clients.

## Why ?

You have a Kafka topic with some consumers and you want to upgrade the Avro schema
for you producer. How do you ensure to not break you consumers with a breaking change
into you schema? Which consumer should you upgrade before?

You need to consume a Kafka topic managed by an another team. How do you ensure
that you current Avro schema version will be able to consume the messages? Which
version should you use? How do you prevent the other team to make a schema change
not handled by your consumer?
