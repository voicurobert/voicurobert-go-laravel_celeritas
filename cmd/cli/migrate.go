package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	switch arg2 {
	case "up":
		return cel.MigrateUp(dsn)
	case "down":
		if arg3 == "all" {
			return cel.MigrateDownAll(dsn)
		} else {
			return cel.Steps(-1, dsn)
		}
	case "reset":
		err := cel.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = cel.MigrateUp(dsn)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}

	return nil
}
