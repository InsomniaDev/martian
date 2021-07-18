import React, { useState } from "react";
import { IconButton, Select, MenuItem, FormControl, InputLabel } from "@material-ui/core";
import { makeStyles } from '@material-ui/core/styles';
import Icon from "@mdi/react";
import { mdiTelevision, mdiTelevisionOff } from "@mdi/js";
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

export function HarmonyComponent({ harmony }) {
    const areaCurrentlyOn = harmony.state !== "-1" ? true : false;

    const [harmonyValue, setHarmonyValue] = React.useState("-1");
    const [active, changeActive] = useState(false);
    const [updateHarmonyActivity] = useMutation(changeDeviceStatus);
    const classes = useStyles();

    const changeHarmonyActivity = (harmonyActivityId) => {
        updateHarmonyActivity({
            variables: {
                id: harmonyActivityId,
                integration: harmony.integration,
                status: harmony.state,
                level: harmonyActivityId
            }
        });
    }

    if (harmony.state !== harmonyValue) {
        setHarmonyValue(harmony.state);
    }

    if (areaCurrentlyOn !== active) {
        changeActive(areaCurrentlyOn);
    }

    const handleHarmonySelection = ({ target }) => {
        setHarmonyValue(target.value)
        changeHarmonyActivity(target.value);
    }

    return (
        <div key={harmony.id + "_harmony_main"} className={classes.controls}>
            <div key={harmony.id + "_harmony_change_button"} onClick={() => {
                if (active) {
                    changeActive(false);
                    return changeHarmonyActivity(0);
                } else {
                    changeActive(true);
                    return changeHarmonyActivity(100);
                }
            }}>
                <IconButton key={harmony.id + "_cardHarmonyButton"} >
                    {active ? (
                        <Icon key={harmony.id + "_cardIcon"} path={mdiTelevision} className={classes.lightOn} size={3} />
                    ) : (
                        <Icon key={harmony.id + "_cardIcon"} path={mdiTelevisionOff} size={3} />
                    )}
                </IconButton>
            </div>
            <form autoComplete='off'>
                <FormControl className={classes.formControl}>
                    <InputLabel >Activity</InputLabel>
                    <Select
                        id="harmony-selection"
                        value={harmonyValue}
                        onChange={handleHarmonySelection}>
                            {
                                JSON.parse(harmony.name).map(a => (
                                    <MenuItem value={a.activityID}>{a.name}</MenuItem>
                                ))
                            }
                    </Select>
                </FormControl>
            </form>
        </div>
    )
}