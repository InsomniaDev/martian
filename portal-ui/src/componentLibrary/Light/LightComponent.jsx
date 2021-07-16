import React, { useState } from "react";
import { IconButton, Typography } from "@material-ui/core";
import { makeStyles } from '@material-ui/core/styles';
import Icon from "@mdi/react";
import { mdiLightbulbOn, mdiLightbulbOff } from "@mdi/js";
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

export function LightComponent({ light }) {
    const areaCurrentlyOn = light.state.toLowerCase() !== "off" ? true : false;

    const [active, changeActive] = useState(false);
    const [changeAreaStatus] = useMutation(changeDeviceStatus);
    const classes = useStyles();

    const changeLightStatus = () => {
        changeAreaStatus({
            variables: {
                id: light.id,
                integration: light.integration,
                status: light.state.toLowerCase() === "on" ? "off" : "on",
                level: ""
            }
        });
    }

    if (areaCurrentlyOn !== active) {
        changeActive(areaCurrentlyOn);
    }

    return (
        <div key={light.name + "_light_main"} className={classes.controls}>
            <div key={light.name + "_light_change_button"} onClick={() => {
                if (active) {
                    changeActive(false);
                    return changeLightStatus();
                } else {
                    changeActive(true);
                    return changeLightStatus();
                }
            }}>
                <IconButton key={light.name + "_cardLightButton"} >
                    {active ? (
                        <Icon key={light.name + "_cardIcon"} path={mdiLightbulbOn} className={classes.lightOn} size={3} />
                    ) : (
                        <Icon key={light.name + "_cardIcon"} path={mdiLightbulbOff} size={3} />
                    )}
                </IconButton>
            </div>
            <Typography key={light.name + "_light_title"} className={classes.title}>
                {light.name}
            </Typography>
        </div>
    )
}