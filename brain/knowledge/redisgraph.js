import { Graph } from "redisgraph.js";
import { Config } from "../../models/config";

const config = new Config().getDatabaseInfo();

export const DeviceGraph = new Graph("Devices", config.url, config.port);
export const EventGraph = new Graph("Events", config.url, config.port);
