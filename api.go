package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type G struct {
	Groups []Group `json:"groups"`
}
type T struct {
	Tasks []Task `json:"tasks"`
}
type CustomTaskPOST struct {
	Title    string `json:"title"`
	Group_id string `json:"group_id"`
}
type CustomTimeframesPOST struct {
	Beginning time.Time `json:"from"`
	Ending    time.Time `json:"to"`
	Task_id   string    `json:"task_id"`
}
type CustomTaskPUT struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Group_id string `json:"group_id"`
}
type Group struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Group_tasks []Task `json:"tasks"`
}

type Task struct {
	ID              int          `json:"id"`
	Title           string       `json:"title"`
	Task_timeframes []Timeframes `json:"time_frames"`
}

type Timeframes struct {
	ID        int       `json:"id"`
	Beginning time.Time `json:"from"`
	Ending    time.Time `json:"to"`
	Task_id   int       `json:"task_id"`
}

func groupsGetHandler(w http.ResponseWriter, r *http.Request) {
	//не учтена возможность одного таска иметь несколько промежутков времени выполнения
	groups := []Group{}
	groupRows, err := db.Query("SELECT * FROM groups LEFT JOIN tasks ON tasks.group_id = groups.id LEFT JOIN timeframes ON timeframes.task_id = tasks.id")
	if err != nil {
		panic(err)
	}
	var counter int = 0
	defer groupRows.Close()
	for groupRows.Next() {
		g := Group{}
		t := Task{}
		tf := Timeframes{}
		tArr := []Task{}
		tfArr := []Timeframes{}
		var temp1 int = 0
		err := groupRows.Scan(&g.ID, &g.Title, &t.ID, &t.Title, &temp1, &tf.ID, &tf.Beginning, &tf.Ending, &tf.Task_id)
		if counter > 0 {
			if g.ID == groups[counter-1].ID {
				tfArr = append(tfArr, tf)
				t.Task_timeframes = tfArr
				tArr = append(tArr, t)
				groups[counter-1].Group_tasks = append(groups[counter-1].Group_tasks, tArr[0])
				counter--
			} else {
				tfArr = append(tfArr, tf)
				t.Task_timeframes = tfArr
				tArr = append(tArr, t)
				g.Group_tasks = tArr
				groups = append(groups, g)
			}
		} else {
			tfArr = append(tfArr, tf)
			t.Task_timeframes = tfArr
			tArr = append(tArr, t)
			g.Group_tasks = tArr
			groups = append(groups, g)
		}
		if err != nil {
			fmt.Println(err)
		}
		counter++
	}
	output := G{}
	output.Groups = groups
	json, err := json.MarshalIndent(output, "\t", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(json))
}

func tasksGetHandler(w http.ResponseWriter, r *http.Request) {
	//не учтена возможность одного таска иметь несколько промежутков времени выполнения
	tasks := []Task{}
	taskRows, err := db.Query("SELECT * FROM tasks LEFT JOIN timeframes ON timeframes.task_id = tasks.id")
	if err != nil {
		panic(err)
	}
	defer taskRows.Close()
	for taskRows.Next() {
		t := Task{}
		tf := Timeframes{}
		tfArr := []Timeframes{}
		var temp1 int = 0
		err := taskRows.Scan(&t.ID, &t.Title, &temp1, &tf.ID, &tf.Beginning, &tf.Ending, &tf.Task_id)
		tfArr = append(tfArr, tf)
		t.Task_timeframes = tfArr
		tasks = append(tasks, t)
		if err != nil {
			fmt.Println(err)
		}
	}
	output := T{}
	output.Tasks = tasks
	json, err := json.MarshalIndent(output, "\t", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(json))
}

func groupsPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var group Group
	_ = json.NewDecoder(r.Body).Decode(&group)
	_, err := db.Exec("INSERT INTO groups(title) VALUES($1)", group.Title)
	if err != nil {
		panic(err)
	}
	groupRows, err := db.Query("SELECT * FROM groups WHERE title = $1", group.Title)
	if err != nil {
		panic(err)
	}
	defer groupRows.Close()
	var output Group
	for groupRows.Next() {
		err = groupRows.Scan(&output.ID, &output.Title)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	json, err := json.MarshalIndent(output, "\t", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(json))
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "status created")
}

func tasksPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task CustomTaskPOST
	_ = json.NewDecoder(r.Body).Decode(&task)
	i, _ := strconv.Atoi(task.Group_id)
	_, err := db.Exec("INSERT INTO tasks(title, group_id) VALUES($1, $2)", task.Title, i)
	if err != nil {
		panic(err)
	}
	taskRows, err := db.Query("SELECT * FROM tasks WHERE title = $1", task.Title)
	if err != nil {
		panic(err)
	}
	defer taskRows.Close()
	var output Task
	var temp int = 0
	for taskRows.Next() {
		err = taskRows.Scan(&output.ID, &output.Title, &temp)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	json, err := json.MarshalIndent(output, "\t", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(json))
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "status created")
}

func timeframesPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var timeframes CustomTimeframesPOST
	_ = json.NewDecoder(r.Body).Decode(&timeframes)
	i, _ := strconv.Atoi(timeframes.Task_id)
	_, err := db.Exec("INSERT INTO timeframes(beginning, ending, task_id) VALUES($1, $2, $3)", timeframes.Beginning, timeframes.Ending, i)
	if err != nil {
		panic(err)
	}
	timeframeRows, err := db.Query("SELECT * FROM timeframes WHERE beginning = $1 AND ending = $2 AND task_id = $3", timeframes.Beginning, timeframes.Ending, i)
	if err != nil {
		panic(err)
	}
	defer timeframeRows.Close()
	var output Timeframes
	for timeframeRows.Next() {
		err = timeframeRows.Scan(&output.ID, &output.Beginning, &output.Ending, &output.Task_id)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	json, err := json.MarshalIndent(output, "\t", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(json))
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, "status created")
}

func groupsPutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	groups := []Group{}
	groupRows, err := db.Query("SELECT * FROM groups")
	if err != nil {
		panic(err)
	}
	defer groupRows.Close()
	for groupRows.Next() {
		g := Group{}
		err := groupRows.Scan(&g.ID, &g.Title)
		if err != nil {
			fmt.Println(err)
			continue
		}
		groups = append(groups, g)
	}
	params := mux.Vars(r)
	group := Group{}
	for i := 0; i < len(groups); i++ {
		id, _ := strconv.Atoi(params["id"])
		if groups[i].ID == id {
			_ = json.NewDecoder(r.Body).Decode(&group)
			group.ID = id
			groups[i] = group
			break
		}
	}
	json.NewEncoder(w).Encode(group)
	_, err = db.Exec("UPDATE groups SET title = $1 WHERE id = $2", group.Title, group.ID)
	if err != nil {
		panic(err)
	}
}

func tasksPutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tasks := []CustomTaskPUT{}
	taskRows, err := db.Query("SELECT * FROM tasks")
	if err != nil {
		panic(err)
	}
	for taskRows.Next() {
		t := CustomTaskPUT{}
		err := taskRows.Scan(&t.ID, &t.Title, &t.Group_id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tasks = append(tasks, t)
	}
	params := mux.Vars(r)
	task := CustomTaskPUT{}
	for i := 0; i < len(tasks); i++ {
		id, _ := strconv.Atoi(params["id"])
		if tasks[i].ID == id {
			_ = json.NewDecoder(r.Body).Decode(&task)
			task.ID = id
			tasks[i] = task
			break
		}
	}
	json.NewEncoder(w).Encode(task)
	_, err = db.Exec("UPDATE tasks SET title = $1, group_id = $2 WHERE id = $3", task.Title, task.Group_id, params["id"])
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, "status ok")
}

func groupsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM groups WHERE id = $1", params["id"])
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, "status 204 no content")
}

func tasksDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM tasks WHERE id = $1", params["id"])
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, "status 204 no content")
}

func timeframesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM timeframes WHERE id = $1", params["id"])
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, "status 204 no content")
}
