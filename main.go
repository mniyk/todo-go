package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

var (
	file_name    string = "todo.json"
	view_flag           = flag.Bool("v", false, "View ToDo.")
	add_flag            = flag.Bool("a", false, "Add ToDo.")
	update_flag         = flag.Bool("u", false, "Update ToDo.")
	delete_flag         = flag.Bool("d", false, "Update ToDo.")
	filter_flag         = flag.Int("f", 0, "View ToDo filter.(0: all, 1: true, 2: false)")
	id_flag             = flag.Int64("i", 0, "ToDo id.")
	content_flag        = flag.String("c", "", "ToDo content.")
	tag_flag            = flag.String("t", "", "ToDo tag.")
	duedate_flag        = flag.String("due", "", "ToDo duedate.(yyyy-mm-dd hh:mm:ss)")
)

type ToDo struct {
	Id        int64     // ID
	Content   string    // 内容
	Done      bool      // 完了フラグ
	Tag       string    // タグ
	DueDate   time.Time // 期日
	UpdatedAt time.Time // 更新日
}

func read_todos(filter int) []ToDo {
	var todos []ToDo

	// jsonファイルの読み込み
	byte_data, err := os.ReadFile(file_name)
	if err != nil {
		// jsonファイルの作成
		file, err := os.Create(file_name)
		if err != nil {
			log.Fatal(err)
		}
		file.Close()
	}

	// jsonデータの構造体インスタンス化
	_ = json.Unmarshal(byte_data, &todos)

	new_todos := []ToDo{}

	for _, v := range todos {
		if filter == 0 {
			new_todos = append(new_todos, v)
		} else if filter == 1 {
			if v.Done {
				new_todos = append(new_todos, v)
			}
		} else if filter == 2 {
			if !v.Done {
				new_todos = append(new_todos, v)
			}
		}
	}

	return new_todos
}

func add_todo(todo ToDo) []ToDo {
	todos := read_todos(0)

	todos = append(todos, todo)

	write_json_data(todos)

	return todos
}

func write_json_data(todos []ToDo) {
	// jsonファイルの取得
	file, _ := os.Create(file_name)
	defer file.Close()

	// jsonファイルへの書き込み
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(todos); err != nil {
		log.Fatal(err)
	}
}

func update_todo(id int64, done bool) []ToDo {
	todos := read_todos(0)
	new_todos := []ToDo{}

	for _, v := range todos {
		if v.Id == id {
			v.Done = done
		}
		new_todos = append(new_todos, v)
	}

	write_json_data(new_todos)

	return new_todos
}

func delete_todo(id int64) []ToDo {
	todos := read_todos(0)
	new_todos := []ToDo{}

	for _, v := range todos {
		if v.Id != id {
			new_todos = append(new_todos, v)
		}
	}

	write_json_data(new_todos)

	return new_todos
}

func get_last_id(todos []ToDo) int64 {
	var last_id int64 = 0

	for _, v := range todos {
		if v.Id >= last_id {
			last_id = v.Id
		}
	}

	return last_id + 1
}

func table_render(todos []ToDo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Content", "Done", "Tag", "DueDate", "UpdatedAt"})

	for _, v := range todos {
		append_data := []string{
			strconv.FormatInt(v.Id, 10),
			v.Content,
			strconv.FormatBool(v.Done),
			v.Tag,
			v.DueDate.Format("2006-01-02 15:04:05"),
			v.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		table.Append(append_data)
	}

	table.Render()
}

func main() {
	flag.Parse()

	if *view_flag {
		todos := read_todos(*filter_flag)

		table_render(todos)
	} else if *add_flag {
		todos := read_todos(0)
		last_id := get_last_id(todos)

		if *content_flag == "" {
			log.Fatalln("No content.")
		}

		if *tag_flag == "" {
			log.Fatalln("No tag.")
		}

		if *duedate_flag == "" {
			log.Fatalln("No duedate.")
		}

		duedate, err := time.Parse("2006-01-02 15:04:05", *duedate_flag)
		if err != nil {
			log.Fatal(err)
		}

		todo := ToDo{
			Id:        last_id,
			Content:   *content_flag,
			Done:      false,
			Tag:       *tag_flag,
			DueDate:   duedate,
			UpdatedAt: time.Now(),
		}

		todos = add_todo(todo)

		table_render(todos)
	} else if *update_flag {
		if *id_flag == 0 {
			log.Fatalln("No ID.")
		}

		todos := update_todo(*id_flag, true)

		table_render(todos)
	} else if *delete_flag {
		if *id_flag == 0 {
			log.Fatalln("No ID.")
		}

		todos := delete_todo(*id_flag)

		table_render(todos)
	}
}
