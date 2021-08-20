import { makeStyles } from '@material-ui/core/styles';
import React from "react";
import classNames from 'classnames';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import Typography from '@material-ui/core/Typography';
import Chip from '@material-ui/core/Chip';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';


const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        padding: "1rem",
        margin: "3%",
        textAlign: 'center',
        backgroundColor: "rgb(199, 216, 223)",
        width: '90%',
    },
    heading: {
        fontSize: theme.typography.pxToRem(23),
    },
    secondaryHeading: {
        fontSize: theme.typography.pxToRem(20),
        color: theme.palette.text.secondary,
    },
    icon: {
        verticalAlign: 'bottom',
        height: 20,
        width: 20,
    },
    details: {
        alignItems: 'center',
    },
    column: {
        flexBasis: '33.33%',
    },
    helper: {
        borderLeft: `2px solid ${theme.palette.divider}`,
        padding: `${theme.spacing()}px ${theme.spacing(2)}px`,
    },
    link: {
        color: theme.palette.primary.main,
        textDecoration: 'none',
        '&:hover': {
            textDecoration: 'underline',
        },
    },
}));



export function IntegrationArea({ integration }) {
    //   const { loading, error, data, refetch,  } = useQuery(getMenuConfiguration, {
    //     pollInterval: 500,
    //     fetchPolicy: "no-cache"
    //   });

    const classes = useStyles();

    // Set up the subscription for the light state
    // subscribeToMore({
    //   document: subscriptionForMenu,
    //   updateQuery: (prev, { subscriptionData: {
    //     data: { menuChange },
    //   } }) => {
    //     if (!menuChange) {
    //       return prev;
    //     }
    //     const newObject = {
    //       "menuConfiguration": menuChange
    //     }
    //     return Object.assign({}, prev, newObject)
    //   }
    // });

    //   if (loading) return <p>loading ...</p>;
    //   if (error) {
    //     console.log(error);
    //     return <p>Error :( {error}</p>;
    //   }

    return (
        <ExpansionPanel >
            <ExpansionPanelSummary >
                <div className={classes.column}>
                    <Typography className={classes.heading}>{integration.name}</Typography>
                </div>
                <div className={classes.column}>
                    <Typography className={classes.secondaryHeading}>Select Devices</Typography>
                </div>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails className={classes.details}>
                <div className={classes.column} />
                <div className={classes.column}>
                    <Chip label="Barbados" className={classes.chip} onDelete={() => { }} />
                </div>
                <div className={classNames(classes.column, classes.helper)}>
                    <Typography variant="caption">
                        Select your destination of choice
                        <br />
                        <a href="#sub-labels-and-columns" className={classes.link}>
                            Learn more
                        </a>
                    </Typography>
                </div>
            </ExpansionPanelDetails>
            <Divider />
            <ExpansionPanelActions>
                <Button size="small">Cancel</Button>
                <Button size="small" color="primary">
                    Save
                </Button>
            </ExpansionPanelActions>
        </ExpansionPanel>
    )
}
