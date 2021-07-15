import React, { useState } from "react";
import { IconButton, Typography } from "@material-ui/core";
import { makeStyles } from '@material-ui/core/styles';
import Icon from "@mdi/react";
import { mdiPowerPlug, mdiPowerPlugOff } from "@mdi/js";
import { useMutation } from "@apollo/client";
import { changeDeviceStatus } from "../mutations/changeDeviceState";

const useStyles = makeStyles((theme) => ({
    controls: {
        display: 'flex',
        alignItems: 'center',
        paddingLeft: theme.spacing(1),
        paddingBottom: theme.spacing(1),
    },
    lightOn: {
        color: "#F8F8FF",
    },
    title: {
        fontSize: '1.4rem',
        color: "#DCDCDC",
        paddingRight: ".5rem",
        marginRight: '.5rem'
    },
}));

export function PlugComponent({ plug }) {
    const areaCurrentlyOn = plug.state.toLowerCase() !== "off" ? true : false;

    const [active, changeActive] = useState(false);
    const [changeAreaStatus] = useMutation(changeDeviceStatus);
    const classes = useStyles();

    const changePlugStatus = () => {
        changeAreaStatus({
            variables: {
                id: plug.id,
                integration: plug.integration,
                status: plug.state.toLowerCase() === "on" ? "off" : "on",
                level: ""
            }
        });
    }

    if (areaCurrentlyOn !== active) {
        changeActive(areaCurrentlyOn);
    }

    return (
        <div key={plug.name + "_plug_main"} className={classes.controls}>
            <div key={plug.name + "_plug_change_button"} onClick={() => {
                if (active) {
                    changeActive(false);
                    return changePlugStatus();
                } else {
                    changeActive(true);
                    return changePlugStatus();
                }
            }}>
                <IconButton key={plug.name + "_cardLightButton"} >
                    {active ? (
                        <Icon key={plug.name + "_cardIcon"} path={mdiPowerPlug} className={classes.lightOn} size={3} />
                    ) : (
                        <Icon key={plug.name + "_cardIcon"} path={mdiPowerPlugOff} size={3} />
                    )}
                </IconButton>
            </div>
            <Typography key={plug.name + "_light_title"} className={classes.title}>
                {plug.name}
            </Typography>
        </div>
    )
}