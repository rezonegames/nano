package main

import (
	"nano"
	"nano/component"
	"nano/examples/gamecluster/db"
	"nano/examples/gamecluster/gate"
	"nano/examples/gamecluster/master"
	"nano/examples/gamecluster/tetris"
	"nano/serialize/json"
	"github.com/pingcap/errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net/http"
	"os"
)

func main()  {
	app := cli.NewApp()
	app.Name = "tetris"
	app.Author = "rezonegames"
	app.Email = "duhe@live.hk"
	app.Description = "Nano cluster demo"
	app.Commands = []cli.Command{
		{
			Name: "master",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen",
					Usage: "Master service listen address",
					Value: "127.0.0.1:34561",
				},
			},
			Action: runMaster,
		},
		{
			Name: "gate",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "master",
					Usage: "master server address",
					Value: "127.0.0.1:34561",
				},
				cli.StringFlag{
					Name:  "listen",
					Usage: "Gate service listen address",
					Value: "127.0.0.1:34567",
				},
				cli.StringFlag{
					Name:  "gate-address",
					Usage: "Client connect address",
					Value: "127.0.0.1:34560",
				},
			},
			Action: runGate,
		},
		{
			Name: "room",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "master",
					Usage: "master server address",
					Value: "127.0.0.1:34561",
				},
				cli.StringFlag{
					Name:  "listen",
					Usage: "Chat service listen address",
					Value: "127.0.0.1:34568",
				},
			},
			Action: runTetris,
		},
	}
	log.SetLevel(log.DebugLevel)
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("Startup server error %+v", err)
	}
}

func runMaster(args *cli.Context) error {
	listen := args.String("listen")
	if listen == "" {
		return errors.Errorf("master listen address cannot empty")
	}
	services := &component.Components{}
	services.Register(master.NewMasterService())
	nano.Listen(listen,
		nano.WithMaster(),
		nano.WithComponents(services),
		//nano.WithDebugMode(),
		nano.WithSerializer(json.NewSerializer()),
	)
	return nil
}

func runGate(args *cli.Context) error {
	listen := args.String("listen")
	if listen == "" {
		return errors.Errorf("master listen address cannot empty")
	}
	gateAddr := args.String("gate-address")
	if gateAddr == "" {
		return errors.Errorf("gate address cannot empty")
	}
	masterAddr := args.String("master")
	if masterAddr == "" {
		return errors.Errorf("master address cannot empty")
	}
	if err := db.NewMongoClient("mongodb://host.docker.internal:27018,host.docker.internal:27017/Client?replicaSet=rs0&readPreference=secondary"); err != nil {
		return errors.Errorf("mongodb init error err v+%", err)
	}

	services := &component.Components{}
	services.Register(gate.NewGateService())

	nano.Listen(listen,
		nano.WithAdvertiseAddr(masterAddr),
		nano.WithClientAddr(gateAddr),
		nano.WithComponents(services),
		nano.WithCheckOriginFunc(func(_ *http.Request) bool { return true }),
		//nano.WithDebugMode(),
		nano.WithSerializer(json.NewSerializer()),
	)
	return nil
}

func runTetris(args *cli.Context) error {
	listen := args.String("listen")
	if listen == "" {
		return errors.Errorf("master listen address cannot empty")
	}
	masterAddr := args.String("master")
	if masterAddr == "" {
		return errors.Errorf("master address cannot empty")
	}
	if err := db.NewMongoClient("mongodb://host.docker.internal:27018,host.docker.internal:27017/Slots5StoreBeta?replicaSet=rs0&readPreference=secondary"); err != nil {
		return errors.Errorf("mongodb init error err v+%", err)
	}

	service := &component.Components{}
	roomservice := tetris.NewRoomService(
		tetris.WithCap(100000),
		tetris.WithTableCap(6),
	)
	tableservice := tetris.NewTableService(
		tetris.WithRoom(roomservice.Room()),
		tetris.WithTableCap(6),
	)
	service.Register(roomservice)
	service.Register(tableservice)

	nano.Listen(listen,
		nano.WithAdvertiseAddr(masterAddr),
		nano.WithComponents(service),
		//nano.WithDebugMode(),
		nano.WithSerializer(json.NewSerializer()),
	)
	return nil
}