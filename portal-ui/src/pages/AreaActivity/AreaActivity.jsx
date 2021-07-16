import { makeStyles } from '@material-ui/core/styles';
import React from "react";
import { useQuery } from "@apollo/client";
import Grid from "@material-ui/core/Grid";
import { Area } from "../../components/Area/Area";
import { subscriptionForMenu } from './subscriptions/menuChangesSubscription';
import { getMenuConfiguration } from './queries/getMenuConfiguration';

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    padding: "1rem",
    margin: "1rem",
    textAlign: 'center',
    backgroundColor: "#02464d",
    fontSize: "2rem",
    display: 'flex',
    flexDirection: 'row'
  },
  paper: {
    padding: theme.spacing(2),
    height: "3rem",
    color: "#F8F8FF",
    backgroundColor: "#008B8B" //00FFFF
  },
  paperChose: {
    padding: theme.spacing(2),
    height: "3rem",
    color: "#F8F8FF",
    backgroundColor: "#9932CC"
  },
  paperAreaActive: {
    padding: theme.spacing(2),
    height: "3rem",
    color: "#F8F8FF",
    backgroundColor: "#9932CC"
  }
}));



export function AreaActivity() {
  const { loading, error, data, refetch } = useQuery(getMenuConfiguration, {
    pollInterval: 500,
    fetchPolicy: "no-cache"
  });

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

  if (loading) return <p>loading ...</p>;
  if (error) {
    console.log(error);
    return <p>Error :( {error}</p>;
  }

  return (
    <div key="main" className={classes.root}>
      <Grid key="main_grid" container spacing={3}>
        {data.menuConfiguration.map(area => (
          <Grid key={area.areaName + "_grid"} item md={4}>
            <Area key={area.areaName} refetch={refetch} area={area}>
            </Area>
          </Grid>
        ))}
      </Grid>
    </div>
  )
}
