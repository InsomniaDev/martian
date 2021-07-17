import React, { useState, useEffect } from "react";
import { makeStyles } from '@material-ui/core/styles';
import { Card, IconButton } from "@material-ui/core";
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import Icon from "@mdi/react";
import { mdiLightbulbOn, mdiLightbulbOff, mdiPowerPlug, mdiPowerPlugOff, mdiFan, mdiFanOff } from "@mdi/js";
import { useMutation } from "@apollo/client";
import { changeDeviceStatus } from "../../componentLibrary/mutations/changeDeviceState";
import ReactCardFlip from 'react-card-flip';
import { LightComponent } from "../../componentLibrary/Light/LightComponent";
import { PlugComponent } from "../../componentLibrary/Plug/PlugComponent";
import { FanComponent } from "../../componentLibrary/Fan/FanComponent";
import AreaMenu from "../AreaMenu/AreaMenu";

const useStyles = makeStyles({
    rootOn: {
        backgroundColor: "rgb(13, 180, 180)",
    },
    root: {
        backgroundColor: "#008B8B" //00FFFF
    },
    lightOn: {
        color: "#F8F8FF",
    },
    title: {
        fontSize: '1.5rem',
        color: "#F8F8FF",
    },
    titleOn: {
        fontSize: '1.5rem',
        color: "rgb(51, 46, 51)",
    },
    backOfCardTitle: {
        fontSize: '1.5rem',
        color: "#F8F8FF",
    },
    pos: {
        marginBottom: 12,
    },
    ellipsis: {
        float: 'right',
        position: 'absolute'
    }
});

export function Area({ refetch, area }) {
    const areaCurrentlyOn = area.active;
    const [active, changeActive] = useState(false);
    const [lightActive, changeLightActive] = useState(false);
    const [plugActive, changePlugActive] = useState(false);
    const [fanActive, changeFanActive] = useState(false);
    const [flipped, changeFace] = useState(false);

    const [lightExists, changeLightExists] = useState(false);
    const [plugExists, changePlugExists] = useState(false);
    const [fanExists, changeFanExists] = useState(false);

    const [changeAreaStatus] = useMutation(changeDeviceStatus);

    if (areaCurrentlyOn !== active) {
        changeActive(areaCurrentlyOn);
    }
    var types = [];
    area.devices.forEach((device) => {
        if (!types.includes(device.type)) {
            types.push(device.type);
        }
    });

    const classes = useStyles();

    useEffect(() => {
        const interval = setInterval((flipped) => {
            changeFace(flipped);
        }, 10000);

        return () => {
            clearInterval(interval);
        };
    });

    var lightOn = false;
    var plugOn = false;
    var fanOn = false;
    area.devices.forEach(device => {
        if (device.type.toLowerCase() === "light") {
            if (lightExists !== true) {
                changeLightExists(true);
            }
            if (device.state.toLowerCase() !== "off") {
                lightOn = true;
            }
        } else if (device.type.toLowerCase() === "plug") {
            if (plugExists !== true) {
                changePlugExists(true);
            }
            if (device.state.toLowerCase() !== "off") {
                plugOn = true;
            }
        } else if (device.type.toLowerCase() === "fan") {
            if (fanExists !== true) {
                changeFanExists(true);
            }
            if (device.state.toLowerCase() !== "off") {
                fanOn = true;
            }
        }
    })
    if (lightOn !== lightActive) {
        changeLightActive(lightOn);
    }
    if (plugOn !== plugActive) {
        changePlugActive(plugOn);
    }
    if (fanOn !== fanActive) {
        changeFanActive(fanOn);
    }

    const bulkChangeAreaStatus = (type) => {
        var typeStatusForArea = false;
        switch (type) {
            case "light":
                typeStatusForArea = lightOn;
                break;
            case "plug":
                typeStatusForArea = plugOn;
                break;
            case "fan":
                typeStatusForArea = fanOn;
                break;
        }
        area.devices.forEach(device => {
            if (device.type.toLowerCase() === type) {
                var status = device.state.toLowerCase() === "off" ? false : true;
                var value = device.value === 100 ? 0 : 100;
                console.log({device})
                console.log({typeStatusForArea})
                if (status === typeStatusForArea) {
                    // Need to set this as a timeout so that we don't write to the websocket at the same time
                    setTimeout(() => changeAreaStatus({
                        variables: {
                            id: device.id,
                            integration: device.integration,
                            status: status ? "off" : "on",
                            level: value === 100 ? 0 : 100,
                        }
                    }), 10);
                }
            }
        })
        refetch();
    }

    const iconSize = 2.5;

    return (
        <ReactCardFlip isFlipped={flipped} flipDirection="vertical">
            <Card className={areaCurrentlyOn ? classes.rootOn : classes.root} key={area.areaName + "_card"} >
                <div onClick={() => changeFace(!flipped)}>
                    <CardContent key={area.areaName + "_cardContent"}>
                        <Typography key={area.areaName + "_cardTitle"} className={areaCurrentlyOn ? classes.titleOn : classes.title} color="textSecondary" gutterBottom>
                            {area.areaName}
                        </Typography>
                    </CardContent>
                </div>
                <CardActions key={area.areaName + "_cardActions"} >
                    {lightExists ? (<div onClick={() => {
                        if (lightActive) {
                            changeLightActive(false);
                            return bulkChangeAreaStatus("light");
                        } else {
                            changeLightActive(true);
                            return bulkChangeAreaStatus("light");
                        }
                    }}>
                        <IconButton key={area.areaName + "_cardLightButton"} >
                            {lightActive ? (
                                <Icon key={area.areaName + "_cardIcon"} path={mdiLightbulbOn} className={classes.lightOn} size={iconSize} />
                            ) : (
                                <Icon key={area.areaName + "_cardIcon"} path={mdiLightbulbOff} size={iconSize} />
                            )}
                        </IconButton>
                    </div>) : (<div></div>)}
                    {plugExists ? (<div onClick={() => {
                        if (plugActive) {
                            changePlugActive(false);
                            return bulkChangeAreaStatus("plug");
                        } else {
                            changePlugActive(true);
                            return bulkChangeAreaStatus("plug");
                        }
                    }}>
                        <IconButton key={area.areaName + "_cardPlugButton"} >
                            {plugActive ? (
                                <Icon key={area.areaName + "_cardIcon"} path={mdiPowerPlug} className={classes.lightOn} size={iconSize} />
                            ) : (
                                <Icon key={area.areaName + "_cardIcon"} path={mdiPowerPlugOff} size={iconSize} />
                            )}
                        </IconButton>
                    </div>) : (<div></div>)}
                    {fanExists ? (<div onClick={() => {
                        if (fanActive) {
                            changeFanActive(false);
                            return bulkChangeAreaStatus("fan");
                        } else {
                            changeFanActive(true);
                            return bulkChangeAreaStatus("fan");
                        }
                    }}>
                        <IconButton key={area.areaName + "_cardFanButton"} >
                            {fanActive ? (
                                <Icon key={area.areaName + "_cardIcon"} path={mdiFan} className={classes.lightOn} size={iconSize} />
                            ) : (
                                <Icon key={area.areaName + "_cardIcon"} path={mdiFanOff} size={iconSize} />
                            )}
                        </IconButton>
                    </div>) : (<div></div>)}
                </CardActions>
            </Card>
            <Card className={classes.root} key={area.areaName + "_card"} >
                <div onClick={() => changeFace(!flipped)}>
                    <CardContent key={area.areaName + "_cardContent"}>
                        <Typography key={area.areaName + "_cardTitle"} className={classes.backOfCardTitle} color="textSecondary" gutterBottom>
                            {area.areaName}
                        </Typography>
                    </CardContent>
                </div>
                {area.devices.map(device => {
                    if (device.type.toLowerCase() === "light") {
                        return (
                            <LightComponent key={device.id} light={device} />
                        )
                    } else if (device.type.toLowerCase() === "plug") {
                        return (
                            <PlugComponent key={device.id} plug={device} />
                        )
                    } else if (device.type.toLowerCase() === "fan") {
                        return (
                            <FanComponent key={device.id} fan={device} />
                        )
                    } else {
                        return (
                            <div></div>
                        )
                    }
                })}
                <AreaMenu area={area} className={classes.ellipsis}></AreaMenu>
            </Card>
        </ReactCardFlip>
    );
}
