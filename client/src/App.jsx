import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import Login from './Login';
import Profile from './Profile';
import RoomList from './RoomList';
import Reservation from './Reservation';

const App = () => {
  return (
    <Router>
      <div className="app">
        <Switch>
          <Route path="/login" component={Login} />
          <Route path="/profile" component={Profile} />
          <Route path="/rooms" component={RoomList} />
          <Route path="/reservation/:roomId" component={Reservation} />
        </Switch>
      </div>
    </Router>
  );
};

export default App;
