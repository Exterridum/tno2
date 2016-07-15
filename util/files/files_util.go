package files

// func loadModels(r *mux.Router) {
// 	files, err := ioutil.ReadDir("models")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var things device.Things
// 	things.Things = make([]string, 0)

// 	for _, fileInfo := range files {
// 		if fileInfo.IsDir() == false {
// 			model := loadModel(fileInfo.Name())

// 			processModel(r, model)

// 			things.Things = append(things.Things, Concat("/", model.ID))
// 		}
// 	}
