package main

import (
  "github.com/go-martini/martini"
  "github.com/codegangsta/martini-contrib/binding"
  "github.com/martini-contrib/render"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
)



// Application refering to a Werker App
type Application struct {
    Id          bson.ObjectId   `bson:"_id"`
    Name        string          `bson:"name"`
    Number      int             `bson:"number"`
}




// DB Returns a martini.Handler
func DB() martini.Handler {
    session, err := mgo.Dial("mongodb://localhost")
    if err != nil {
        panic(err)
    }

    return func(c martini.Context) {
        s := session.Clone()
        c.Map(s.DB("buildnumber"))
        defer s.Close()
        c.Next()
    }
}



func main() {
  m := martini.Classic()
  
  m.Use(render.Renderer())
  m.Use(DB())


  m.Get("/", func() string {
    return "Welcome to build number server";
  })
  m.Get("/api", func(r render.Render) {
     r.JSON(200, map[string]interface{}{"message": "No direct access"})
  })



  // Lists applications
  m.Get("/api/applications", func(r render.Render, db *mgo.Database) {
    var apps []Application
    db.C("applications").Find(nil).All(&apps)
    r.JSON(200, apps)
  })


  //Creates application
  m.Post("/api/applications", 
    binding.Bind(Application{}), 
    func(r render.Render, db *mgo.Database, app Application) {
      //app.Number = 0
      db.C("applications").Insert(&app)
      r.JSON(201, app)
  })


   //Creates new build on application
  m.Post("/api/applications/:id/inc", 
    func(r render.Render, db *mgo.Database, params martini.Params) {
    
   


    var app Application
    
    change := mgo.Change{
        Update: bson.M{"$inc": bson.M{"number": 1}},
        ReturnNew: true,
    }
    
    db.C("applications").FindId(bson.ObjectIdHex(params["id"])).Apply(change,&app)  




    //return db.C("applications").Find(bson.M{"_id": "543f0ba9cdaa667d250a62db"}).Count();
    r.JSON(200, app)
  })



  m.Run()

   
}