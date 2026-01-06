
import { Button, Card, Input, List } from 'antd';

const data = [
  {
    title: '站点名称',
  },
  {
    title: '查单ip',
  },
  {
    title: '代理流量ip',
  },
  {
    title: '代理套餐ip',
  },
  {
    title: '支付超时时间(秒)',
  },
  {
    title: '查单超时时间(秒)',
  },
  {
    title: '50',
  },
  {
    title: '100',
  },
  {
    title: '200',
  },
  {
    title: '300',
  },
  {
    title: '500',
  },
  {
    title: '600',
  },
  {
    title: '700',
  },
  {
    title: '800',
  },
  {
    title: '900',
  },
  {
    title: '1000',
  },
  {
    title: 'ip数量限制',
  },
  {
    title: '小号轮询',
  },
  {
    title: '缓存订单过期时间',
  },
  {
    title: '小号省份配置',
  },
];

function SystemConfig() {
  // const [list, setList] = useState<DataType[]>([])

  // const fetchListRealNameAccount = async () => {
  //   const { data } = await listRealNameAccount()
  //   let d: DataType[] = data?.list?.map((item, index) => {
  //     let newItem: DataType = {
  //       key: index.toString(),
  //       ...item
  //     }
  //     return newItem
  //   })
  //   setList(d)
  // }

  // useEffect(() => {
  //   fetchListRealNameAccount()
  // }, [])

  return (
    <>
      <Card>
        <List
          size="small"
          itemLayout="horizontal"
          dataSource={data}
          renderItem={(item) => (
            <List.Item>
              <span style={{ width: '20%' }}>{item.title}</span>
              <span style={{ width: '60%' }}>
                <Input value={item.title} disabled />
              </span>
              <span style={{ display: 'flex', width: '20%', justifyContent: 'center', alignItems: 'center'}}>
                <Button block type="primary" htmlType="submit" style={{ width: 100 }}>
                  修改
                </Button>
              </span>
            </List.Item>
          )}
        />
      </Card>
    </>
  )
}

export default SystemConfig