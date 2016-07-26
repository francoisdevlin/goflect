Please use the godoc tool to find the documetation for this project.  Start with the main documentation


Things to generate

* Tools for generating type info
  * Type info factories
    * Go type parser
    * Type JSON Marshall/UnMarshall (DONE!  Thanks Go!)
  * Type info sanity
    * Combine with unit test, verify annotations are correct
    * Only one autoinc
    * Nominal with unique field
  * Type conventions
    * Primary Key convention (Id)
    * Nominal Convention

* Consumers of type info
  * CLI interface (Think Action Framework)
  * Curses Interface
  * GUI
    * Web Component
    * Ext Model
    * Ext Store
    * Ext Table
    * Ext Default Form
    * Ext Nominal Dropdown
    * Ext Nominal Autocomplete
  * Marshalling
    * JSON Marshall/UnMarshall (leverage existing, quick hack)
    * XML Marshall/UnMarshall (leverage existing, quick hack)
    * Protobuff Marshall/UnMarshall
    * Avro Marshall/UnMarshall
    * EDN Marshall/UnMarshall
    * CSV Marshall/UnMarshall (leverage existing encoding/csv)
  * SQL
    * Table creation
    * Insert
    * READ 1
    * READ WHERE (Excel autofilter)
    * Update 1
    * Delete 1
  * Standard REST endpoints
    * Create
    * Read
    * Update
    * Delete
    * List Where...
    * Live Docs
    * Nominal Read endpoint
