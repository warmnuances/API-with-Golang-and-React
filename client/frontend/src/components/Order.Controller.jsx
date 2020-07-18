import React from 'react'
import { Segment, Search, Dropdown} from 'semantic-ui-react'
import { withRouter } from 'react-router-dom'
import { useDebounce } from 'use-debounce'
import Axios from "axios";
import OrderTable from './OrderTable'
import _ from 'lodash'
import Skeleton from 'react-loading-skeleton';
import { Grid, Icon, Pagination, Input, Button } from 'semantic-ui-react'
import { DateRange } from 'react-date-range';
import { format, addMonths, getMonth } from 'date-fns';



const url = "http://localhost:8080/orders/all"


function OrderController(props) {
  const { location, history } = props

  const [noPages,setNoPages] = React.useState(1)
  const [search, setSearch] = React.useState("")
  const [loading, setLoading] = React.useState(false)
  const [debouncedSearch] = useDebounce(search, 500);
  const [results, setResult] = React.useState([])


  /**Pages**/  
  const [page,setPage] = React.useState(1)


  /** Datepicker**/
  const [date, setDate] = React.useState([
    {
      startDate: null,
      endDate: null,
      key: 'selection'
    }
  ])

  /** Field**/
  const [field, setField] = React.useState("order_name")
  const [amtField, setAmtField] = React.useState("total")

  /** Range Slider**/
  const [amount,setAmount] = React.useState({
    max: 0,
    min: 0
  })

  const [render, reRender] = React.useState(false)
 
  

  const fieldOptions = [
    {key: 'order_name', value: 'order_name',text: 'Order Name'},
    {key: 'customer_id', value: 'cust_id',text: 'Customer Name'}
  ]

  const amtFieldOptions = [
    {key: 'delivered', value: 'delivered',text: 'Delivered Amount'},
    {key: 'total', value: 'total',text: 'Total Amount'},
  ]


  React.useEffect(() => {
    reRender(false)
    const source = Axios.CancelToken.source();
    setLoading(true)
    getOrders(source.token, page, date, amount, field, debouncedSearch, amtField)
      .then(({data}) => {
        setResult(data.result)
        setLoading(false)
        setNoPages(Math.ceil(data.rowCount/5))
      })
      .catch(e => {
        if (Axios.isCancel(source)) {
          return;
        }
        setResult([]);
      
      });

    return () => {
      setLoading(false)
      source.cancel(
        "Canceled because of component unmounted or debounce Text changed"
      );
    };

  }, [debouncedSearch, page, render, noPages])



  function getOrders(token,page, ...args){
    const {startDate, endDate} = args[0][0]
    const {min, max} = args[1]
    const field = args[2]
    const search = args[3]
    const amtField = args[4]

    
    let sDate = new Date(startDate)
    let eDate = new Date(endDate)
    const formattedStart = format(sDate, "yyyy-mm-dd'T'HH:mm:'12Z'")
    const formattedEnd = format(eDate, "yyyy-mm-dd'T'HH:mm:'12Z'")

    //JS month start with 0 -> normalise
    let nsDate = normaliseDate(formattedStart)
    let neDate = normaliseDate(formattedEnd)

    let formedUrl = `${url}?page=${(page - 1)||0}`
      + `${ (min > 0)? `&${amtField}_min=${min}` : ''}`
      + `${ (max > 0)? `&${amtField}_max=${max}`: ''}`
      + `${(field && search)? `&${field}=${search}`: ""}`
      + `${(startDate && endDate)? `&startDate=${nsDate}&endDate=${neDate}`: ""}`
      

    return Axios
      .get(encodeURI(formedUrl), {
        cancelToken: token
      })
  }

  const handleSearchChange = (e, { value }) => {
    setLoading(true)
    setSearch(value)
  }
 
  function onAmountSubmit(e, item){
    if(amount.min > amount.max ){
      setAmount({
        min: 0,
        max: 0
      })
      alert("Invalid input!")
      return
    }else{
      reRender(true)
    }
  }
  function onAmountChange(e,item){
    setAmount({
      ...amount,
      [e.target.name]: item.value
    })
  }

  function normaliseDate(date){
    let replace = (parseInt(date[6]) + 1) + ""
    let result  = date.substring(0, 6) + replace +date.substring(7)
    return result
  }


  function onClickRender(){
    reRender(true)
  }
  return (
    <Segment>
      <Grid>
        <Grid.Row className="container__filter">
          <div className="container__filter--left"> 
            <h6 className="container__filter--text">Search: </h6>
            <div className="container__filter__element--search">
              <Dropdown
                placeholder='Select fields'
                selection
                defaultValue="order_name"
                options={fieldOptions}
                onChange={(e, {value}) => {
                  setField(value)
                }}
              />
                <Search
                    open={false}
                    loading={loading}
                    onSearchChange={_.debounce(handleSearchChange, 300 ,{leading: true})}
                    value={search}
                    placeholder="Search"
                >   
              </Search>
              <Button primary onClick={onClickRender}>Search</Button>
            </div>
            
            <Grid.Row className="container__filter__element--range">
              <Grid.Row>
                <h5>Amount: </h5>
              </Grid.Row>
              <Dropdown
                  placeholder='Select fields'
                  selection
                  defaultValue="total"
                  options={amtFieldOptions}
                  onChange={(e, {value}) => {
                    setAmtField(value)
                  }}
                />
              <Input placeholder='Minimum...' 
                label="Min Amount"
                value={amount.min} 
                onChange={onAmountChange}
                onKeyDown={(evt) => ["e", "E", "+", "-"].includes(evt.key) && evt.preventDefault()}
                name="min" 
                step={10}
                type="number"/>
              <Input placeholder='Maximum...' 
                label="Max Amount"
                value={amount.max} 
                onChange={onAmountChange}
                onKeyDown={(evt) => ["e", "E", "+", "-"].includes(evt.key) && evt.preventDefault()}
                name="max" 
                step={10}
                type="number"/>
              <Button primary onClick={onAmountSubmit}>Search</Button>
            </Grid.Row>
          </div>
           

          <div className="container__filter--right">
            <h6 className="container__filter--text">Date:</h6>
            <DateRange
                editableDateInputs={true}
                onChange={item => setDate([item.selection])}
                moveRangeOnFirstSelection={false}
                ranges={date}
              />
              <Button primary onClick={onClickRender}>Select Date</Button>
            </div>
        </Grid.Row>

      </Grid>

      {
        loading? 
        <>
          <br/>
          <Skeleton count={5}/> 
        </>:
        <OrderTable data={results} setLoading={setLoading}/>
      }

      <Pagination
          defaultActivePage={page}
          ellipsisItem={{ content: <Icon name='ellipsis horizontal' />, icon: true }}
          firstItem={{ content: <Icon name='angle double left' />, icon: true }}
          lastItem={{ content: <Icon name='angle double right' />, icon: true }}
          prevItem={{ content: <Icon name='angle left' />, icon: true }}
          nextItem={{ content: <Icon name='angle right' />, icon: true }}
          totalPages={noPages}
          onPageChange={(e,{activePage}) => {
            setLoading(true)
            history.push({
              pathname: location.pathname,
              search: `?page=${activePage}`
            })
            setPage(activePage)
          }}
        />
      

    </Segment>
  )
}

export default withRouter(OrderController)
