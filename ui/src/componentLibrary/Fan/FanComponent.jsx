import React, { useState } from "react";
import { IconButton, Typography, Select, MenuItem, FormControl, InputLabel } from "@material-ui/core";
import { makeStyles } from '@material-ui/core/styles';
import Icon from "@mdi/react";
import { mdiFan, mdiFanOff } from "@mdi/js";
import { useMutation } from "@apollo/client";
import { changeDeviceStatus } from "../mutations/changeDeviceState";

const useStyles = makeStyles((theme) => ({
    controls: {
        display: 'flex',
        alignItems: 'center',
        paddingLeft: theme.spacing(1),
        paddingBottom: theme.spacing(1),
        flexWrap: 'wrap',
    },
    rowFlex: {
        flexDirection: 'row',
    },
    lightOn: {
        color: "#F8F8FF",
    },
    formControl: {
        minWidth: 150,
    },
    title: {
        fontSize: '1.4rem',
        color: "#DCDCDC",
        paddingRight: ".5rem",
        marginRight: '.5rem',
    },
}));

export function FanComponent({ fan }) {
    const areaCurrentlyOn = fan.state.toLowerCase() !== "off" ? true : false;

    const [fanValue, setFanValue] = React.useState(0);
    const [active, changeActive] = useState(false);
    const [updateTheFanSTate] = useMutation(changeDeviceStatus);
    const classes = useStyles();

    const changeFanStatus = (level) => {
        var fanStatus = "dim";
        if (level === 0 || level === 100) {
            fanStatus = level === 0 ? "off" : "on";
        }
        updateTheFanSTate({
            variables: {
                id: fan.id,
                integration: fan.integration,
                status: fanStatus,
                level: level
            }
        });
    }

    if (areaCurrentlyOn !== active) {
        changeActive(areaCurrentlyOn);
    }

    if (fan.value !== fanValue) {
        setFanValue(fan.value);
    }

    const handleFanSelection = ({ target }) => {
        setFanValue(target.value)
        changeFanStatus(target.value);
    }

    return (
        <div key={fan.name + "_fan_main"} className={classes.controls}>
            <div key={fan.name + "_fan_change_button"} onClick={() => {
                if (active) {
                    changeActive(false);
                    return changeFanStatus(0);
                } else {
                    changeActive(true);
                    return changeFanStatus(100);
                }
            }}>
                <IconButton key={fan.name + "_cardFanButton"} >
                    {active ? (
                        <Icon key={fan.name + "_cardIcon"} path={mdiFan} className={classes.lightOn} size={3} />
                    ) : (
                        <Icon key={fan.name + "_cardIcon"} path={mdiFanOff} size={3} />
                    )}
                </IconButton>
            </div>
            <Typography key={fan.name + "_fan_title"} className={classes.title}>
                {fan.name}
            </Typography>
            <form autoComplete='off'>
                <FormControl className={classes.formControl}>
                    <InputLabel >Speed</InputLabel>
                    <Select
                        id="fan-selection"
                        value={fanValue}
                        onChange={handleFanSelection}>
                        <MenuItem value={0}>OFF</MenuItem>
                        <MenuItem value={25}>LOW</MenuItem>
                        <MenuItem value={50}>MEDIUM</MenuItem>
                        <MenuItem value={75}>MEDIUM-HIGH</MenuItem>
                        <MenuItem value={100}>HIGH</MenuItem>
                    </Select>
                </FormControl>
            </form>
        </div>
    )
}