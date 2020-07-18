import React from 'react';
import { Route, Switch, Link } from 'react-router-dom';
import { Menu } from 'semantic-ui-react'
import 'semantic-ui-css/semantic.min.css'

 // main style file
import 'react-date-range/dist/styles.css';
import 'react-date-range/dist/theme/default.css'


//Components
import OrderController from './components/Order.Controller'
import HomePage from './components/HomePage'

function App() {
  return (
    <section className="App">
      <Header />
      <Switch>
        <Route exact path="/" component={HomePage}/>
        <Route exact path="/Orders" component={OrderController}/>
      </Switch>
    </section>
  );
}

function Header(){
  const [active, setActive] = React.useState("Home")
  const handleItemClick = (e, { name }) => setActive(name)

  return (
    <Menu pointing secondary>
      <Menu.Item
        as={Link}
        to="/"
        name='Home'
        active={active === 'Home'}
        content='Home'
        onClick={handleItemClick}
      />

      <Menu.Item
        name='Orders'
        as={Link}
        to="/Orders"
        active={active === 'Orders'}
        content='Orders'
        onClick={handleItemClick}
      />
    </Menu>
  )
}


export default App;
