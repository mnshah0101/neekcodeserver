package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateSolution(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateSolution: Received request")

	question_id := r.FormValue("question_id")
	if question_id == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		log.Println("CreateSolution: Missing question_id")
		return
	}

	username := r.FormValue("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		log.Println("CreateSolution: Missing username")
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		log.Println("CreateSolution: Missing email")
		return
	}

	date := r.FormValue("date")
	if date == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		log.Println("CreateSolution: Missing date")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		log.Println("CreateSolution: Error reading file -", err)
		return
	}
	defer file.Close()
	log.Println("CreateSolution: File uploaded with name", header.Filename)

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		log.Println("CreateSolution: Error reading file content -", err)
		return
	}

	// Upload file to S3
	fileName := header.Filename
	fileURL, err := UploadFileToS3(fileName, fileContent)
	if err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		log.Println("CreateSolution: Error uploading file to S3 -", err)
		return
	}
	log.Println("CreateSolution: File uploaded to S3 at", fileURL)

	// Save solution to MongoDB
	collection := MongoClient.Database("test").Collection("solutions")
	solution := Solution{
		ID:          primitive.NewObjectID(),
		Username:    username,
		Email:       email,
		File_url:    fileURL,
		Likes:       0,
		Question_id: question_id,
		Date:        date,
	}

	_, err = collection.InsertOne(context.Background(), solution)
	if err != nil {
		http.Error(w, "Error saving solution", http.StatusInternalServerError)
		log.Println("CreateSolution: Error inserting solution into MongoDB -", err)
		return
	}
	log.Println("CreateSolution: Solution saved to MongoDB")

	response := map[string]interface{}{
		"message": "Solution created successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Println("CreateSolution: Response sent successfully")
}

func GetSolutions(w http.ResponseWriter, r *http.Request) {
	log.Println("GetSolutions: Received request")

	question_id := r.URL.Query().Get("question_id")
	if question_id == "" {
		http.Error(w, "Question ID is required", http.StatusBadRequest)
		log.Println("GetSolutions: Missing question_id")
		return
	}

	collection := MongoClient.Database("test").Collection("solutions")
	cursor, err := collection.Find(context.Background(), bson.M{"question_id": question_id})
	if err != nil {
		http.Error(w, "Error fetching solutions", http.StatusInternalServerError)
		log.Println("GetSolutions: Error fetching solutions from MongoDB -", err)
		return
	}
	defer cursor.Close(context.Background())

	var solutions []Solution
	for cursor.Next(context.Background()) {
		var solution Solution
		if err := cursor.Decode(&solution); err != nil {
			log.Println("GetSolutions: Error decoding solution -", err)
			continue
		}
		solutions = append(solutions, solution)
	}
	log.Println("GetSolutions: Solutions fetched successfully")

	response := map[string]interface{}{
		"solutions": solutions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Println("GetSolutions: Response sent successfully")
}

func LikeSolution(w http.ResponseWriter, r *http.Request) {
	log.Println("LikeSolution: Received request")

	solution_id := r.FormValue("solution_id")
	if solution_id == "" {
		http.Error(w, "Solution ID is required", http.StatusBadRequest)
		log.Println("LikeSolution: Missing solution_id")
		return
	}

	collection := MongoClient.Database("test").Collection("solutions")
	solutionID, err := primitive.ObjectIDFromHex(solution_id)
	if err != nil {
		http.Error(w, "Invalid solution ID format", http.StatusBadRequest)
		log.Println("LikeSolution: Invalid solution ID format -", err)
		return
	}

	filter := bson.M{"_id": solutionID}
	update := bson.M{"$inc": bson.M{"likes": 1}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error liking solution", http.StatusInternalServerError)
		log.Println("LikeSolution: Error updating likes in MongoDB -", err)
		return
	}
	log.Println("LikeSolution: Solution liked successfully")

	response := map[string]interface{}{
		"message": "Solution liked successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Println("LikeSolution: Response sent successfully")
}

func DislikeSolution(w http.ResponseWriter, r *http.Request) {
	log.Println("DislikeSolution: Received request")

	solution_id := r.FormValue("solution_id")
	if solution_id == "" {
		http.Error(w, "Solution ID is required", http.StatusBadRequest)
		log.Println("DislikeSolution: Missing solution_id")
		return
	}

	collection := MongoClient.Database("test").Collection("solutions")
	solutionID, err := primitive.ObjectIDFromHex(solution_id)
	if err != nil {
		http.Error(w, "Invalid solution ID format", http.StatusBadRequest)
		log.Println("DislikeSolution: Invalid solution ID format -", err)
		return
	}

	filter := bson.M{"_id": solutionID}
	update := bson.M{"$inc": bson.M{"likes": -1}}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, "Error disliking solution", http.StatusInternalServerError)
		log.Println("DislikeSolution: Error updating dislikes in MongoDB -", err)
		return
	}
	log.Println("DislikeSolution: Solution disliked successfully")

	response := map[string]interface{}{
		"message": "Solution disliked successfully",
		"status":  "success",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	log.Println("DislikeSolution: Response sent successfully")
}
