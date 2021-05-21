package main

import (
    "database/sql"
    "fmt"
    "strconv"
    "log"
    _ "github.com/mattn/go-sqlite3"
)

func getHostDetails(hostName string) map[string]string {
        var allVars = make(map[string]string)
        //set local variables
        var hostVarName string
        var hostVarValue string
       //processing query result
       database, _ := sql.Open("sqlite3", "./ansible.db")
       hostVars, err :=
	       database.Query("SELECT host_var, value FROM hostvars WHERE host_id = (SELECT rowid FROM hosts WHERE name='" + hostName + "');")
	if err != nil {
		fmt.Println("SQL select error: ")
		log.Fatal(err)
	}
	for hostVars.Next() {
		hostVars.Scan(&hostVarName, &hostVarValue)
		allVars[hostVarName] = hostVarValue
	}
        return allVars
}

func getHostList() map[string]interface{} {
        var allHost []string
	var hostName string
        database, _ := sql.Open("sqlite3", "./ansible.db")
        fetchAllHost, err :=
                database.Query("SELECT name FROM hosts ORDER BY name ASC")
        if err != nil {
		fmt.Println("SQL select error: ")
		log.Fatal(err)
	}
        for fetchAllHost.Next() {
                fetchAllHost.Scan(&hostName)
                allHost = append(allHost,hostName)
        }
        hostMap := map[string]interface{}{
                "hosts": allHost,
        }
	return hostMap
}

func getMetaData() map[string]interface{} {
	//var hostMetaData = make(map[string]interface{})
	var hostMap = make(map[string]interface{})
	var allVars = make(map[string]string)
        database, _ := sql.Open("sqlite3", "./ansible.db")
        allHost, err :=
                database.Query("SELECT rowid,name FROM hosts ORDER BY name ASC")
        if err != nil {
                fmt.Println("SQL select error: ")
                log.Fatal(err)
        }
        for allHost.Next() {
		//set local variables
		var hostVarName string
		var hostVarValue string
		var hostId int
		var hostName string
		//processing query result
		allHost.Scan(&hostId,&hostName)
		hostVars, err :=
                database.Query("SELECT host_var, value FROM hostvars WHERE host_id=" + strconv.Itoa(hostId))
	        if err != nil {
			fmt.Println("SQL select error: ")
			log.Fatal(err)
		}
		for hostVars.Next() {
			hostVars.Scan(&hostVarName, &hostVarValue)
                        allVars[hostVarName] = hostVarValue
                }
                hostMap[hostName] = allVars
	}
	return hostMap
}

func getAllList() map[string]interface{} {
//Set local variables
        var parentName string
        var groupName string
        var parentVarName string
        var parentVarValue string
	var combinedMap = make(map[string]interface{})
//Set all host list
        combinedMap["all"] = getHostList()
// Set all host metadata
        combinedMap["_meta"] = getMetaData()
//Fetch all Parent groups and handle SQL error
        database, _ := sql.Open("sqlite3", "./ansible.db")
        rows, err :=
                database.Query("SELECT rowid,name FROM parentgroups ORDER BY name ASC")
        if err != nil {
                fmt.Println("SQL select error: ")
                log.Fatal(err)
        }
//Create map from all children of the Parent
        for rows.Next() {
	        var allGroup []string
	        var allVars = make(map[string]string)
		var rowId int
		rows.Scan(&rowId, &parentName)
		groupAssoc, err :=
			database.Query("SELECT name FROM hostgroups WHERE rowid IN (SELECT group_id FROM parent_child WHERE parent_id=" + strconv.Itoa(rowId) + ");")
	        if err != nil {
			fmt.Println("SQL select error: ")
			log.Fatal(err)
		}
		for groupAssoc.Next() {
			groupAssoc.Scan(&groupName)
			allGroup = append(allGroup,groupName)
		}
		parentVars, err :=
			database.Query("SELECT parent_var, value FROM parentvars WHERE parent_id=" + strconv.Itoa(rowId))
		if err != nil {
			fmt.Println("SQL select error: ")
			log.Fatal(err)
		}
		for parentVars.Next() {
			parentVars.Scan(&parentVarName, &parentVarValue)
			allVars[parentVarName] = parentVarValue
		}
		parentMap := map[string]interface{}{
			"children": allGroup,
			"vars": allVars,
                }
		combinedMap[parentName] = parentMap
	}
        return combinedMap
}

