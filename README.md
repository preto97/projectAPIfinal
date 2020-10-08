
# Project APIfinal

  In this Project I created restful API to capture and display the body of any HTTP POST request made to the app.
  The API is linked to a database, named "testAPI", created in mySQL.
  // Into the testAPI.sql DB, I created a table named "songs" with fields: - id
  //                                                                       - title
  //                                                                       - duration
  //                                                                       - singer
  // Before you run the app make sure you activate mySQL service from xampp.
  // The appliation start a localhost:8080 which can accept different patterns that makes specific actions to table "songs":
  // Some actions can run from the web browser and the others one, runs just from the Postman app.
  // So, to the localhost:8080 you can add the following patterns (this runs from web browser because have GET method request):   
  //  - "/songs"          -> to see all songs stored in the tabel "songs" 
  //  - "/songs/{id}"     -> to see a specific song, by the id
  //  - "/songss/delAll"  -> to delete all songs stored in table "songs"
  // Attention! All the next patterns doesn't have GET method request, like above. So, you can run them just from Postman:
  //  - "/songs"        -> (request method: POST)   -> you can add a song to the table "songs", inserting to the body fields, your datas
  //  - "/songs/{id}"   -> (request method: PUT)    -> you can update an existing record from the table "songs", by the id
  //                                                   The update action, executes after you had insert to the body fields, your new datas
  //  - "/songs/{id}"   -> (request method: DELETE) -> you can delete an existing record from the table "songs", by the id
  //
  // All the fields from table "songs" are encoded to json, so the content-type header is set to "application/json".
  //
  // In the future I want to extend this project, I want to implement all actions that executs in Postman to specific buttons into the web browser.
  // Also I want to change the update action, beacuse for now it has PUT request and I want to change it to PATCH, for keep 
  // the old datas, if I don't change all fields in the update body.
  // As well, I want to convert/change the "duration" field type from string to time.
  
           
  
  
