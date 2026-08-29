package router

import (
	"encoding/json"
	"log"
	"os"

	"github.com/LeBaoTai/myco-controller/internal/controller/gnmib"
	"github.com/LeBaoTai/myco-controller/internal/oc"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/ygot/ytypes"
	"go.yaml.in/yaml/v4"
)

type Session struct {
	conn *gnmib.Client
}

type ConfigPath struct {
	Paths []string `yaml:"paths"`
}

func CreateNewSession(conn *gnmib.Client) *Session {
	return &Session{
		conn: conn,
	}
}

func getConfigPath() []*gnmi.Path {
	pathFile, err := os.ReadFile("./internal/config/paths.yml")
	if err != nil {
		log.Panic(err)
	}

	var configPath ConfigPath
	err = yaml.Unmarshal(pathFile, &configPath)
	if err != nil {
		log.Fatalf("Cannot parse the yml file: %v", err)
	}

	var pathList []*gnmi.Path
	for _, path := range configPath.Paths {
		builtPath, err := gnmib.BuildPath(path)
		if err != nil {
			log.Printf("Cannot create path: %v\n", err)
		}
		pathList = append(pathList, builtPath)
	}
	return pathList
}

func (s *Session) HandleChageConfiguration(update []byte) {
	// get configuration path from config yaml file
	res, err := s.conn.GetConfig(getConfigPath())
	if err != nil {
		log.Printf("Cannot get config: %v\b", err)
	}

	// creat root current currentSchema and go struct device
	currentSchema, err := oc.Schema()
	if err != nil {
		log.Printf("Cannot get schema: %v\n", err)
	}
	//currentSystem := &oc.Device{}
	//currentSchema.Root = currentSystem

	// get current config
	for _, n := range res.Notification {
		err := ytypes.UnmarshalNotifications(
			currentSchema,
			[]*gnmi.Notification{n},
			&ytypes.IgnoreExtraFields{},
		)
		if err != nil {
			log.Println(err)
		}
	}

	// create incoming change
	// newSystem := &oc.Device{}
	// if err := oc.Unmarshal(update, newSystem); err != nil {
	// 	log.Printf("Cannot parse the new config: %v\n", err)
	// }
	//
	// diff, err := ygot.Diff(currentSystem, newSystem)
	// if err != nil {
	// 	log.Printf("Cannot compare the config: %v", err)
	// }

	// log.Printf("Config Changes: %v", diff)

	// err = ytypes.UnmarshalNotifications(
	// 	currentSchema,
	// 	res.Notification,
	// 	&ytypes.IgnoreExtraFields{},
	// )
	// if err != nil {
	// 	log.Printf("Failed to unmarshal path data: %v", err)
	// }
	//
	// jsonOpts := &ygot.EmitJSONConfig{
	// 	Format: ygot.RFC7951,
	// 	Indent: "  ",
	// }
	//
	// jsonString, err := ygot.EmitJSON(currentSystem, jsonOpts)
	// if err != nil {
	// 	log.Printf("Failed:%v", err)
	// }
}

func convertGetResponseToYAML(res *gnmi.GetResponse) ([]byte, error) {
	var sb []byte
	for _, notif := range res.Notification {
		for _, upd := range notif.Update {
			raw := upd.Val.GetJsonIetfVal()
			var data any
			if err := json.Unmarshal(raw, &data); err != nil {
				return []byte(""), err
			}
			out, err := yaml.Marshal(data)
			if err != nil {
				return []byte(""), err
			}
			sb = append(sb, out...)
		}
	}
	return sb, nil
}

func (s *Session) GerCurrentConfiguration() ([]byte, error) {
	res, err := s.conn.GetConfigAsASCII(getConfigPath())
	if err != nil {
		log.Printf("Cannot get config: %v\b", err)
	}

	// convert gnmi to byte
	var sb []byte
	for _, notif := range res.Notification {
		for _, upd := range notif.Update {
			raw := upd.Val.GetJsonIetfVal()
			var data any
			if err := json.Unmarshal(raw, &data); err != nil {
				return []byte(""), err
			}
			out, err := yaml.Marshal(data)
			if err != nil {
				return []byte(""), err
			}
			sb = append(sb, out...)
		}
	}

	return sb, nil
}

func (s *Session) Close() {
	s.conn.Close()
}
