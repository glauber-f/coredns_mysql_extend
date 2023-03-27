package coredns_mysql_extend

import (
	"strconv"
	"time"

	"github.com/coredns/caddy"
)

var (
	zoneQuerySQL   string
	recordQuerySQL string
)

func (m *Mysql) Name() string {
	return pluginName
}
func (m *Mysql) parseConfig(c *caddy.Controller) error {
	mysqlConfig := mysqlConfig{
		dsn:          defaultDSN,
		dumpFile:     defaultDumpFile,
		ttl:          defaultTTL,
		zonesTable:   defaultZonesTable,
		recordsTable: defaultRecordsTable,

		maxIdleConns:    defaultMaxIdleConns,
		maxOpenConns:    defaultMaxOpenConns,
		connMaxIdleTime: defaultConnMaxIdleTime,
		connMaxLifetime: defaultConnMaxLifeTime,

		failHeartbeatTime:    defaultFailHeartBeatTime,
		successHeartbeatTime: defaultSuccessHeartBeatTime,
	}

	m.mysqlConfig = mysqlConfig
	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "dsn":
				if !c.NextArg() {
					return c.ArgErr()
				}
				m.dsn = c.Val()
			case "dump_file":
				if !c.NextArg() {
					return c.ArgErr()
				}
				m.dumpFile = c.Val()
			case "ttl":
				if !c.NextArg() {
					return c.ArgErr()
				}
				userTTL, err := strconv.Atoi(c.Val())
				if err != nil || userTTL <= zero {
					m.ttl = defaultTTL
				} else {
					m.ttl = uint32(userTTL)
				}
			case "zones_table":
				if !c.NextArg() {
					return c.ArgErr()
				}
				m.zonesTable = c.Val()
			case "records_table":
				if !c.NextArg() {
					return c.ArgErr()
				}
				m.recordsTable = c.Val()
			case "db_max_idle_conns":
				if !c.NextArg() {
					return c.ArgErr()
				}
				userMaxIdleConns, err := strconv.Atoi(c.Val())
				if err != nil || userMaxIdleConns <= zero {
					m.maxIdleConns = defaultMaxIdleConns
				} else {
					m.maxIdleConns = userMaxIdleConns
				}
			case "db_max_open_conns":
				userMaxOpenConns, err := strconv.Atoi(c.Val())
				if err != nil || userMaxOpenConns <= zero {
					m.maxOpenConns = defaultMaxOpenConns
				} else {
					m.maxOpenConns = userMaxOpenConns
				}
			case "db_conn_max_idle_time":
				userConnMaxIdleTime, err := time.ParseDuration(c.Val())
				if err != nil || userConnMaxIdleTime <= zeroTime {
					m.connMaxIdleTime = defaultConnMaxIdleTime
				} else {
					m.connMaxIdleTime = userConnMaxIdleTime
				}
			case "db_conn_max_life_time":
				userConnMaxLifeTime, err := time.ParseDuration(c.Val())
				if err != nil || userConnMaxLifeTime <= zeroTime {
					m.connMaxLifetime = defaultConnMaxLifeTime
				} else {
					m.connMaxLifetime = userConnMaxLifeTime

				}
			case "fail_heartbeat_time":
				userFailHeartBeatTime, err := time.ParseDuration(c.Val())
				if err != nil || userFailHeartBeatTime <= zeroTime {
					m.failHeartbeatTime = defaultFailHeartBeatTime
				} else {
					m.failHeartbeatTime = userFailHeartBeatTime
				}
			case "success_heartbeat_time":
				userSuccessHeartBeatTime, err := time.ParseDuration(c.Val())
				if err != nil || userSuccessHeartBeatTime <= zeroTime {
					m.successHeartbeatTime = defaultSuccessHeartBeatTime
				} else {
					m.successHeartbeatTime = userSuccessHeartBeatTime
				}
			default:
				return c.Errf("unknown property '%s'", c.Val())
			}
		}
	}
	return nil
}

func (m *Mysql) createTables() error {
	_, err := m.DB.Exec(`
        CREATE TABLE IF NOT EXISTS ` + m.zonesTable + ` (
            id INT NOT NULL AUTO_INCREMENT,
            zone_name VARCHAR(255) NOT NULL,
            PRIMARY KEY (id),
            UNIQUE KEY (zone_name)
        );
    `)
	if err != nil {
		logger.Error(err)
		return err
	}

	_, err = m.DB.Exec(`
        CREATE TABLE IF NOT EXISTS ` + m.recordsTable + ` (
            id INT NOT NULL AUTO_INCREMENT,
            zone_id INT NOT NULL,
            hostname VARCHAR(512) NOT NULL,
            type VARCHAR(10) NOT NULL,
            data VARCHAR(1024) NOT NULL,
            ttl INT NOT NULL DEFAULT 120,
			online INT NOT NULL DEFAULT 0,
            PRIMARY KEY (id),
            FOREIGN KEY (zone_id) REFERENCES ` + m.zonesTable + `(id)
        );
    `)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}