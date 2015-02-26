/*
This is the package goflect, which is used to abstract out many common pattern in application development using reflection.

Overview

There are several common problems assocaited with writing an application, which lead to people writing frameworks.  At SevOne, our applications have several requirements that must be met over and over again.  Some of the most interesting ones include

	* Marshalling data from an external data source
	* Validating that the data is well formed
	* Combining ad hoc user queries with permissions requirements
	* Intercepting, filtering, and routing streams of data
	* Dispatching data across a cluster
	* Providing consistent API endpoint
	* Providing the ability to inspect a process to determine how to interact with it
	* Providing consistant machine readable data structures to determine how applications interact

Marshalling Data

One of the first issues is providing a consistent way to convert binary data at the perimiter of an application to a normal Go type.  This typically leads to a lot of boilerplate code to do this.  You can use the sql package to get information out of the database using reflection.  We can use the build in json & xml packages to handle encoding those formats

Validating Data with Relational Algebra

Once the records are marshalled, ew need to ensure that they are well formed.  Ensuring that the data is well formed is one of the biggest security concerns in modern applications.  It is very easy to specify most constraints using relational algebra, but most applications don't.  The goflect package is designed to provide a powerful relational algebra tool that can be used to add further constraints on data.  These constraints are composible, enabling them to be reused at both design time and run time

Ad hoc user requests

Users like being able to make add hoc requests, which is great.  Relational algebra is composable, so it should be very easy to chain any requests that users have for our data with any addition constraints in our system.  However, ANSI SQL is far from composable, and prevent many security risks.  Again, the matcher package helps us solve this problem

Authorization concerns

Authorization is a conern about who is allowed to access what records.  It involves creating an additional algebraic constraint, at run time.  These constaints are complete orthogonal to the algebraic constraints provided by both the data restrictions and user requests.  Once again, the matcher package can be used to provide an algebra to cleanly describe authorization

Stream filtering and routing

Here comes the broken record.  We have a constraint that we would like to apply to a stream, and filter/route/multiplex appropriately.  The matcher API can be used to describe the rules that govern these streams, and then route the packages appropriately

Dispatching data across a cluster

Dispatching data across a cluster is a common problem that we have to resolve at SevOne.  Does a record get inserted in to the master, the local box, or a specific peer?  Do we connect directly to a mysql instance, or need to go through an api?

These problems can be solved by providing a common record manipulation interface that clients consume.  This interface will be database technology, rest back, location agnostic.  This is possible through two things:

	* Migrating all relational algebra to the matcher package
	* Providing a record service API that consumes records and matchers, and does not tell about anything else.

The ability to swap out back end services on a whim, and abstract out dispatch is huge.  This will help streamline testing, and development.  This will also allow our product to play in a software defined data center out of the box

Providing consistent API endpoints

Once we move all dispatch to a central API, we can then have consumers stop worrying about implementations.  This will allow API endpoints such as REST and SOAP to be applied automatically on to of the central API.  This will also provide a new level of consistency to our API client

Provide process documentation

The matcher API comes with a pretty printer.  This will allow us to explain to consumers WHY the system is behaving the way it is.  Along with using the record definition extensions, we can then generate lot of very useful human facing artifacts:

	* An Admin GUI
	* A "Bulk Change" tool, like JIRA
	* A CLI
	* Endpoint documentation

Also, it will be possible to produce a machine readable version of the field documentation as well.  If all of our applications were to provide the same endpoint documentation, we could then automatically have a map of how all of the system processess interact

*/
package main

import (
//"database/sql"
//"strings"
//"html/template"
//"net/http"
//"flag"
//"fmt"
//"os"
//"reflect"
//"strconv"
)
