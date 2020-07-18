import React from 'react'
import { Table } from 'semantic-ui-react'
import { format } from 'date-fns';

const dateFormat = "MMM dd, h:mm:a"

function Orders(props) {
  const { data, setLoading } = props; 
  // console.log(data)
  React.useEffect(() => {
    setLoading(false)
  },[])

  return (
    <Table singleLine>
      <Table.Header>
        <Table.Row>
          <Table.HeaderCell>Order Name</Table.HeaderCell>
          <Table.HeaderCell>Customer Company</Table.HeaderCell>
          <Table.HeaderCell>Customer Name</Table.HeaderCell>
          <Table.HeaderCell>Order Date</Table.HeaderCell>
          <Table.HeaderCell>Delivered Amount</Table.HeaderCell>
          <Table.HeaderCell>Total Amount</Table.HeaderCell>
        </Table.Row>
      </Table.Header>

      <Table.Body>
        {
          (data.length > 0)?
            data.map((item,idx) => {
              return(
                <Table.Row key={idx}>
                  <Table.Cell>{item.order_name}</Table.Cell>
                  <Table.Cell>{item.customer_company}</Table.Cell>
                  <Table.Cell>{item.customer_id}</Table.Cell>
                  <Table.Cell>{format(new Date(item.created_at), dateFormat)}</Table.Cell>
                  <Table.Cell>$ {item.delivered_amount}</Table.Cell>
                  <Table.Cell>$ {item.total_amount}</Table.Cell>
                </Table.Row>
              )
            }):
            <h1> There is no data</h1>
        }
        
      </Table.Body>
    </Table>
  )
}

export default React.memo(Orders)
