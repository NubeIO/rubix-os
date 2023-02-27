# CHANGELOG

## [v0.10.1](https://github.com/NubeIO/flow-framework/tree/v0.10.1) (2023-02-27)

- Remove point write buffers
  - They need to be straightforwardly written & as they are not frequently called too
- Remove fromPlugin flag on PointWrite
- Remove the afterRealDeviceUpdate flag, it's making things confusing

## [v0.10.0](https://github.com/NubeIO/flow-framework/tree/v0.10.0) (2023-02-24)

- Fix: runtime issue
- Fix: parallel sync network and device
- Attach default args of SynNetworkDevices & SyncDevicePoints
- fixes NaN crash on transformations
- Buffer implementation on UpdatePoint and PointWrite
- Improvements on automapping logs
- Have a common UpdatePoint function to send fromPlugin flag
  - galvintmv plugin was attaching fromPlugin=false

## [v0.9.35](https://github.com/NubeIO/flow-framework/tree/v0.9.35) (2023-02-17)

- Use point HistoryEnable, HistoryType and HistoryInterval on producer (#866)
- update to default sch freq
- updates to galvintmv plugin for commissioning
- Fix: MQTT connection reset by peer issue

## [v0.9.34](https://github.com/NubeIO/flow-framework/tree/v0.9.34) (2023-02-16)

- bacnet master: adds priority 7 as write
- adds running and fault properties to most plugins

## [v0.9.33](https://github.com/NubeIO/flow-framework/tree/v0.9.33) (2023-02-13)

- Sync meta-tags and tags

## [v0.9.32](https://github.com/NubeIO/flow-framework/tree/v0.9.32) (2023-02-10)

- support for known lorawan string payloads

## [v0.9.31](https://github.com/NubeIO/flow-framework/tree/v0.9.31) (2023-02-08)

- Adds rtc timezone point
- Updates rtc timezone offset
- Fix: database lock issue

## [v0.9.30](https://github.com/NubeIO/flow-framework/tree/v0.9.30) (2023-02-07)

- Modifies BACnet Master point write strategy

## [v0.9.29](https://github.com/NubeIO/flow-framework/tree/v0.9.29) (2023-02-06)

- Applies schedule timezone to the timestamps
- Fix: concurrent map writes
- Clean up of bacnet master
- Don't build unused plugins
- Build edge28 & rubixio for only armv7

## [v0.9.28](https://github.com/NubeIO/flow-framework/tree/v0.9.28) (2023-01-30)

- Add mqtt connect retry and auto reconnect
- Fix: MQTT subscribers listening

## [v0.9.27](https://github.com/NubeIO/flow-framework/tree/v0.9.27) (2023-01-30)

- Updated bacnetmaster to use all priorities
- Updates to schedules

## [v0.9.26](https://github.com/NubeIO/flow-framework/tree/v0.9.26) (2023-01-28)

- Added mqtt api for schedules

## [v0.9.25](https://github.com/NubeIO/flow-framework/tree/v0.9.25) (2023-01-25)

- Add meta tags query param on network, device and point #824
- Lora updates v1 #825
- small updates to schedule #822
- update to bacnet master #821

## [v0.9.24](https://github.com/NubeIO/flow-framework/tree/v0.9.24) (2023-01-19)

- update modbus schema #819
- Add flow_https_local,flow_ip_local and flow_port_local #818
- added api's for mapping for all protocals #815
- added api for bacnet whois #814

## [v0.9.23](https://github.com/NubeIO/flow-framework/tree/v0.9.23) (2023-01-16)

- Fix: FN selection doesn't exist issue
- Upgrade lib-schema to v0.1.8
- Get writer by name

## [v0.9.22](https://github.com/NubeIO/flow-framework/tree/v0.9.22) (2023-01-16)

- add point write message ok
- Improvements: add waits on loop delete
- Improvements: add waits on loraraw add device points
- added api for bacnet whois
- Add an auto_mapping example for lora and system plugins
- adds safety to modbus byte conversions

## [v0.9.21](https://github.com/NubeIO/flow-framework/tree/v0.9.21) (2023-01-05)

- Improvement: MQTT points list update only on point name changes
- Add auto-mapping
- Sync details on the producer
- Fix: Present value and priority array discrepancies
- Remove the magic string from configs

## [v0.9.20](https://github.com/NubeIO/flow-framework/tree/v0.9.20) (2022-12-20)

- bump version

## [v0.9.19](https://github.com/NubeIO/flow-framework/tree/v0.9.19) (2022-12-20)

- mqtt support get selected points

## [v0.9.18](https://github.com/NubeIO/flow-framework/tree/v0.9.18) (2022-12-12)

- Remove suffix slash (/) from APIs for to support reverse proxy
- prevents clearing of errors on poll queue re-add

## [v0.9.17](https://github.com/NubeIO/flow-framework/tree/v0.9.17) (2022-12-07)

- Make it installable on old distribution

## [v0.9.16](https://github.com/NubeIO/flow-framework/tree/v0.9.16) (2022-12-07)

- fixes modbus polling order.
- adds history analysis plugins (as examples)

## [v0.9.15](https://github.com/NubeIO/flow-framework/tree/v0.9.15) (2022-12-02)

- adds support for modbus bitwise.
- fixed up MaxPollRate on polling

## [v0.9.14](https://github.com/NubeIO/flow-framework/tree/v0.9.14) (2022-11-29)

- networklinker fix name overwrite #772

## [v0.9.13](https://github.com/NubeIO/flow-framework/tree/v0.9.13) (2022-11-25)

- networklinker allow edit point name #770

## [v0.9.12](https://github.com/NubeIO/flow-framework/tree/v0.9.12) (2022-11-23)

- fix edgeinflux cov
- fix config file reset to default
- improve history api query

## [v0.9.11](https://github.com/NubeIO/flow-framework/tree/v0.9.11) (2022-11-21)

- fix edgeinflux history

## [v0.9.10](https://github.com/NubeIO/flow-framework/tree/v0.9.10) (2022-11-18)

- adds more modbus baud rates

## [v0.9.9](https://github.com/NubeIO/flow-framework/tree/v0.9.9) (2022-11-16)

- adds network history enable

## [v0.9.8](https://github.com/NubeIO/flow-framework/tree/v0.9.8) (2022-11-16)

- updates network options on plugins
- adds edgeazure (for testing)

## [v0.9.7](https://github.com/NubeIO/flow-framework/tree/v0.9.7) (2022-11-09)

- slows lorawan and galvintmv api calls
- lorawan plugin gets cs token at runtime

## [v0.9.6](https://github.com/NubeIO/flow-framework/tree/v0.9.6) (2022-11-08)

- makes galvin setup steps run in order

## [v0.9.5](https://github.com/NubeIO/flow-framework/tree/v0.9.5) (2022-11-08)

- fix nil errors on lorawan

## [v0.9.4](https://github.com/NubeIO/flow-framework/tree/v0.9.4) (2022-11-08)

- Typo fix on galvintmv plugin

## [v0.9.3](https://github.com/NubeIO/flow-framework/tree/v0.9.3) (2022-11-07)

- Improves maplora and mapmodbus plugins

## [v0.9.2](https://github.com/NubeIO/flow-framework/tree/v0.9.2) (2022-11-07)

- Improves rubixpointsync plugin

## [v0.9.1](https://github.com/NubeIO/flow-framework/tree/v0.9.1) (2022-11-07)

- Improves polling (bacnet and modbus)
- Improves galvintmv plugin

## [v0.9.0](https://github.com/NubeIO/flow-framework/tree/v0.9.0) (2022-11-07)

- Improves polling (bacnet and modbus)

## [v0.8.9](https://github.com/NubeIO/flow-framework/tree/v0.8.9) (2022-11-04)

- Fixes cancel of MQTT subscribe
- Fixes Rubix Legacy API Calls

## [v0.8.8](https://github.com/NubeIO/flow-framework/tree/v0.8.8) (2022-11-04)

- Re adds legacy mapping plugins

## [v0.8.7](https://github.com/NubeIO/flow-framework/tree/v0.8.7) (2022-11-03)

- Removes legacy mapping plugins

## [v0.8.6](https://github.com/NubeIO/flow-framework/tree/v0.8.6) (2022-11-03)

- Try rebuild to fix modbus plugin (issue since 0.8.4)
- Fixes edgeinflux plugin nil pointer

## [v0.8.5](https://github.com/NubeIO/flow-framework/tree/v0.8.5) (2022-11-02)

- Review unique constraints (#707)

## [v0.8.4](https://github.com/NubeIO/flow-framework/tree/v0.8.4) (2022-11-02)

- Adds plugins to convert legacy lora and modbus to FF

## [v0.8.3](https://github.com/NubeIO/flow-framework/tree/v0.8.3) (2022-10-31)

- Fix: get the latest history for postgres sync

## [v0.8.2](https://github.com/NubeIO/flow-framework/tree/v0.8.2) (2022-10-31)

- added mapping apis (will be removed)

## [v0.8.1](https://github.com/NubeIO/flow-framework/tree/v0.8.1) (2022-10-30)

- Fix: postgres large data push and history mismatched order
- Fix: galvintmv for network search
- Fix: bug on mqtt publish cov

## [v0.8.0](https://github.com/NubeIO/flow-framework/tree/v0.8.0) (2022-10-27)

- adds mqtt to plugins
- adds apis for histories

## [v0.7.9](https://github.com/NubeIO/flow-framework/tree/v0.7.9) (2022-10-25)

- capitalizes lora naming
- adds mqtt functionality

## [v0.7.8](https://github.com/NubeIO/flow-framework/tree/v0.7.8) (2022-10-17)

- adds rubix sync plugin

## [v0.7.7](https://github.com/NubeIO/flow-framework/tree/v0.7.7) (2022-10-10)

- zip-hydro-tap lora fix #671
- plugin/networklinker #670
- updates to edgeinflux, galvintmv, and lorawan #668

## [v0.7.6](https://github.com/NubeIO/flow-framework/tree/v0.7.5) (2022-10-06)

- minor gavlintmv plugin update

## [v0.7.5](https://github.com/NubeIO/flow-framework/tree/v0.7.5) (2022-10-05)

- adds gavlintmv plugin
- adds edgeinflux plugin

## [v0.7.4](https://github.com/NubeIO/flow-framework/tree/v0.7.4) (2022-09-12)

- fix voltage error #652
- fix bug on is mqtt connected #651
- lorawan recurse nested maps in uplink #649

## [v0.7.3](https://github.com/NubeIO/flow-framework/tree/v0.7.3) (2022-09-05)

- bacnetserver repost point values on server restart
- bacnetmaster add multi-state object type support

## [v0.7.2](https://github.com/NubeIO/flow-framework/tree/v0.7.2) (2022-08-31)

- fixed nil pointers and looping in pollqueue
- updated lora raw defaults and payload structures

## [v0.7.1](https://github.com/NubeIO/flow-framework/tree/v0.7.1) (2022-08-24)

- added bacnet multistate

## [v0.7.0](https://github.com/NubeIO/flow-framework/tree/v0.7.0) (2022-08-22)

- added restart plugin api

## [v0.6.9](https://github.com/NubeIO/flow-framework/tree/v0.6.9) (2022-08-22)

- fix bug on bacnet polling

## [v0.6.8](https://github.com/NubeIO/flow-framework/tree/v0.6.8) (2022-08-22)

- updates to all the schemas
- fix on bacnet-master polling
- setup of bacnet-server

## [v0.6.7](https://github.com/NubeIO/flow-framework/tree/v0.6.7) (2022-08-18)

- allow user to delete network if plugin is not installed
- got bacnet-server plugin working with the new bacnet-server app
- remove backup plugin
- Merge pull request #606 from NubeIO/cleanup-of-mapping
- add bacnet read priority
- Merge pull request #608 from NubeIO/list-serial-ports
- Merge pull request #610 from NubeIO/redo-json-schema
- Merge pull request #611 from NubeIO/resync-bacnet-names
- Merge pull request #612 from NubeIO/fix-bacnet-master-net
- Merge pull request #613 from NubeIO/bump-schema
- Fix: remove GlobalUUID unique constraint
- Fix: remove GlobalUUID unique constraint
- Fix: remove consumer's producer_uuid unique index

## [v0.6.6](https://github.com/NubeIO/flow-framework/tree/v0.6.6) (2022-08-09)

- Issue/producer history current writer UUID #598
- Improvements/misc #596
- Add central history producer enable flag #594

## [v0.6.5](https://github.com/NubeIO/flow-framework/tree/v0.6.5) (2022-07-29)

- poll queue nil guarding

## [v0.6.4](https://github.com/NubeIO/flow-framework/tree/v0.6.4) (2022-07-29)

- Improvements to poll queue, modbus, and bacnet

## [v0.6.3](https://github.com/NubeIO/flow-framework/tree/v0.6.3) (2022-07-28)

- Fixed bacnet polling
- Implemented point value transformations (factor, scale, limit, offset)

## [v0.6.2](https://github.com/NubeIO/flow-framework/tree/v0.6.2) (2022-07-27)

- Merge pull request #579 from NubeIO/add/ff-error-drill-down-functions
- Merge pull request #580 from NubeIO/update/seperate-database-methods-from-pollqueue

## [v0.6.1](https://github.com/NubeIO/flow-framework/tree/v0.6.1) (2022-07-19)

- Modbus write improved

## [v0.6.0](https://github.com/NubeIO/flow-framework/tree/v0.6.0) (2022-07-14)

- Merge pull request #528 from NubeIO/redo-lorawan
- added api for polling stats
- Improvements: update history logs after histories gets stored
- Fix: conflict issue on bulk create
- Fix: writeValue and point.WriteValue discrepancy resulting infinite looping (#566)
- fixed edge28 and point update functions

## [v0.5.9](https://github.com/NubeIO/flow-framework/tree/v0.5.9) (2022-07-12)

- improvements to edge28 plugin stability
- improvements to point write for priority array modes
- improvements to modbus plugin
- added bulk interval histories
- improved CreateInBatches
- added interface errors to bacnet network
- added new JSON schema for network, device, point

## [v0.5.8](https://github.com/NubeIO/flow-framework/tree/v0.5.8) (2022-07-05)

- fix on PG DB resync
- clean up of bacnet master

## [v0.5.7](https://github.com/NubeIO/flow-framework/tree/v0.5.7) (2022-07-05)

- added history api to get histories by producer name
- updated point self-mapping to allow enable of history on producer

## [v0.5.6](https://github.com/NubeIO/flow-framework/tree/v0.5.6) (2022-06-17)

- minor modbus updates

## [v0.5.5](https://github.com/NubeIO/flow-framework/tree/v0.5.5) (2022-06-10)

- redo of bacnet-server plugin
- fix bug on lora plugin
- reformat nil helpers
- add in sync on flow-networks

## [v0.5.4](https://github.com/NubeIO/flow-framework/tree/v0.5.4) (2022-05-24)

- bacnet master fix bug on select object type

## [v0.5.3](https://github.com/NubeIO/flow-framework/tree/v0.5.3) (2022-05-24)

- cascade delete on flow networks
- added native bacnet master plugin

## [v0.5.2](https://github.com/NubeIO/flow-framework/tree/v0.5.2) (2022-05-19)

- Make configurable modbus log level

## [v0.5.1](https://github.com/NubeIO/flow-framework/tree/v0.5.1) (2022-05-19)

- Minor update to polling (safeties on timeouts)
- Added a Postgres plugin

## [v0.5.0](https://github.com/NubeIO/flow-framework/tree/v0.5.0) (2022-05-17)

- fixed modbus point/device/network time settings
- updated pointWrite() for modbus compatibility
- priority array utilities added

## [v0.4.9](https://github.com/NubeIO/flow-framework/tree/v0.4.9) (2022-05-13)

- Fix: isChange checker for COV updates

## [v0.4.8](https://github.com/NubeIO/flow-framework/tree/v0.4.8) (2022-05-12)

- Control COV streams on point value update

## [v0.4.7](https://github.com/NubeIO/flow-framework/tree/v0.4.7) (2022-05-12)

- update to modbus-server plugin

## [v0.4.6](https://github.com/NubeIO/flow-framework/tree/v0.4.6) (2022-05-11)

- added new plugin modbus server
- nullable fallback point field support
- bug fixes on points for priority array
- polling for modbus
- bulk write api on points

## [v0.4.5](https://github.com/NubeIO/flow-framework/tree/v0.4.5) (2022-05-02)

- get bacnet auto mapping working to rubix-io
- fix bug on rubix-io plugin
- fix bug on flow-network sync

## [v0.4.4](https://github.com/NubeIO/flow-framework/tree/v0.4.4) (2022-04-29)

- add new plugin for the rubix-io

## [v0.4.3](https://github.com/NubeIO/flow-framework/tree/v0.4.3) (2022-04-26)

- many updates since last build :)

## [v0.4.2](https://github.com/NubeIO/flow-framework/tree/v0.4.2) (2022-03-23)

- Plugin/testing #417
- Improvement/database plugins #415
- small updates to system plugin #414
- added point math func on point write value #413
- updates to modbus, lora and bacnet #412
- updates to lora and bacnet #411

## [v0.4.1](https://github.com/NubeIO/flow-framework/tree/v0.4.1) (2022-03-17)

- Feature: add present_value on writer #405
- small updates to lora #404
- Improvement/db plugins #403
- Improvement/db plugins #402
- remove on device model AddressUUID unique #401
- Improvement/history #399
- Improvement: nullable address_uuid #397
- fix get networks #396
- Support OR query on client_id, site_id & device_id #395

## [v0.4.0](https://github.com/NubeIO/flow-framework/tree/v0.4.0) (2022-02-28)

- Supporting older schedule deployment

## [v0.3.9](https://github.com/NubeIO/flow-framework/tree/v0.3.9) (2022-02-28)

- Added back GET/PATCH by point's name

## [v0.3.8](https://github.com/NubeIO/flow-framework/tree/v0.3.8) (2022-02-24)

- Improvements/misc #385

## [v0.3.7](https://github.com/NubeIO/flow-framework/tree/v0.3.7) (2022-02-24)

- Improvements/misc #384
- Improvement/stream sync #383
- updates to modbus #382
- Improvement/point value update #381
- Replace modbus lib #379

## [v0.3.6](https://github.com/NubeIO/flow-framework/tree/v0.3.6) (2022-02-22)

- Merge pull request #358 from NubeIO/bac-master-1
- Fix: edit flow-network issue
- Fix: issue on P2P openvpn connected devices
  - SyncFlowNetwork: Post "10.8.1.1:1616/ff/api/sync/flow_network";: EOF
- updates to modbus plugin
- Add writers write, read, sync action
- Fix: writer action support for old deployments

## [v0.3.5](https://github.com/NubeIO/flow-framework/tree/v0.3.5) (2022-02-16)

- Improvement: support multiple producers on a single point
- Improvement on producer history
- Writer action updates its own side thing
- Make schedule write actions working
- Fix: FlowNetwork creation issue for HTTP
- Make FlowNetwork update working
- Fix: float pointer values comparison for COV

## [v0.3.4](https://github.com/NubeIO/flow-framework/tree/v0.3.4) (2022-02-10)

- Sync value on COV of point
- Fix: point.present_value comparison issue
- Fix: database lock issue

## [v0.3.3](https://github.com/NubeIO/flow-framework/tree/v0.3.3) (2022-02-10)

- Fix: schedule

## [v0.3.2](https://github.com/NubeIO/flow-framework/tree/v0.3.2) (2022-02-10)

- small bug fix to schedule small fix to stop sch crashing #352

## [v0.3.1](https://github.com/NubeIO/flow-framework/tree/v0.3.1) (2022-02-10)

- small bug fix to schedule small fix to stop sch crashing #351

## [v0.3.0](https://github.com/NubeIO/flow-framework/tree/v0.3.0) (2022-02-10)

- Add schedule writer POC
- Add with_priority option on device query builder
- Improvements on gorm migration #347
- added bacnetmaster plugin #346
- make sure if pnt is same the addrID is not same #345
- clean up of bacnet-server #343
- Update/schedule checker to new schedule json schema #342
- small fixes to modbus #337
- Improvements on schedule #334
- Update to mqtt broker plugin
- Add schedule config
- Sync schedules values on the producer side
- Remove CreateWriterWizard
- Datastore is nil for updating writers (datastore update is only from write actions)
- Add sync on patch (#326)
- Marc/edge28 plugin scaling (#314)
- Merge pull request #325 from NubeIO/improvement/return-appropriate-status-code

## [v0.2.2](https://github.com/NubeIO/flow-framework/tree/v0.2.2) (2021-12-23)

- Improvement on schedule APIs
- Add scheduler for refreshing token

## [v0.2.1](https://github.com/NubeIO/flow-framework/tree/v0.2.1) (2021-12-18)

- Update: update to schedules

## [v0.2.0](https://github.com/NubeIO/flow-framework/tree/v0.2.0) (2021-12-18)

- Update: update to schedules

## [v0.1.9](https://github.com/NubeIO/flow-framework/tree/v0.1.9) (2021-12-16)

- Remove: rubix plugins

## [v0.1.8](https://github.com/NubeIO/flow-framework/tree/v0.1.71) (2021-11-23)

- Fix: get config on sessions

## [v0.1.71](https://github.com/NubeIO/flow-framework/tree/v0.1.71) (2021-11-22)

- made rubix a network #304
- Breaking issue fix on modubs polling #302
- Improvements/misc #301
- Improvement/misc #300
- Add history influx log #299
- Fix: NubeIO vs NubeDev package #298
- Feature/schema api #297
- Update go.mod #296
- updated the bash script #295
- Improvement/misc #294
- Close DB connection #293
- sample api helper #290

## [v0.1.6](https://github.com/NubeIO/flow-framework/tree/v0.1.6) (2021-11-04)

- add APIs for proxying fn, fnc
- serial port fix (#288)

## [v0.1.5](https://github.com/NubeIO/flow-framework/tree/v0.1.5) (2021-11-01)

- nubeio-rubix-lib-helpers-go version upgrade to v0.1.2

## [v0.1.4](https://github.com/NubeIO/flow-framework/tree/v0.1.4) (2021-10-25)

- rubix plugin build fix

## [v0.1.3](https://github.com/NubeIO/flow-framework/tree/v0.1.3) (2021-10-25)

- added rubix-service api
- fix up on schedules
- improvement on writer & writer_clone args query

## [v0.1.2](https://github.com/NubeIO/flow-framework/tree/v0.1.2) (2021-10-21)

- added flow network mqtt api

## [v0.1.1](https://github.com/NubeIO/flow-framework/tree/v0.1.1) (2021-10-19)

- added api for milo db
- clean up of bacnetserver plugin
- added system and time api
- added schedules api

## [v0.1.0](https://github.com/NubeIO/flow-framework/tree/v0.1.0) (2021-10-12)

- fix issues on droplet motion
- added writer-action as thingClass schedule

## [v0.0.9](https://github.com/NubeIO/flow-framework/tree/v0.0.9) (2021-10-11)

- updates to lora and modbus plugins
- added edge-28 plugin

## [v0.0.8](https://github.com/NubeIO/flow-framework/tree/v0.0.8) (2021-10-08)

- updates to lora and modbus plugins

## [v0.0.7](https://github.com/NubeIO/flow-framework/tree/v0.0.7) (2021-10-05)

- updates to lora and modbus plugin

## [v0.0.6](https://github.com/NubeIO/flow-framework/tree/v0.0.6) (2021-10-05)

- fix bug on action write

## [v0.0.5](https://github.com/NubeIO/flow-framework/tree/v0.0.5) (2021-10-05)

- added point calc's, units, and eval
- clean up on lora and modbus plugins

## [v0.0.4](https://github.com/NubeIO/flow-framework/tree/v0.0.4) (2021-10-02)

- added git plugin and updates to modbus

## [v0.0.3](https://github.com/NubeIO/flow-framework/tree/v0.0.3) (2021-10-01)

- make that artifacts working for armv7

## [v0.0.2](https://github.com/NubeIO/flow-framework/tree/v0.0.2) (2021-09-29)

- include plugins on artifacts

## [v0.0.1](https://github.com/NubeIO/flow-framework/tree/v0.0.1) (2021-08-26)

- first initial release
