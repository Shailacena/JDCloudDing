import { Typography, Divider, Row, Col, Card, Statistic, StatisticProps } from 'antd';
import { useAppContext } from '../AppProvider';
import CountUp from 'react-countup';
import { useApis } from '../api/api';
import { useEffect, useState } from 'react';
import { GetOrderSummaryResp } from '../api/types';
const { Title } = Typography;

function Dashboard() {
  let ctx = useAppContext();
  let { getOrderSummary } = useApis()
  const [summary, setSummary] = useState<GetOrderSummaryResp>({});
  const [isLoading, setLoading] = useState<boolean>(false);

  useEffect(() => {
    fetcOrderSummary()
  }, [])

  const fetcOrderSummary = async () => {
    setLoading(true)
    try {
      const { data } = await getOrderSummary()

      setSummary(data)
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false)
    }
  }

  const formatter: StatisticProps['formatter'] = (value) => (
    <CountUp end={value as number} separator="," />
  );

  return (
    <>
      <Typography>
        <Title>欢迎</Title>
        <Title type="danger" level={4}>{ctx?.cookie?.nickname}</Title>
        {/* <Paragraph>当前时间： {new Date().toString()}</Paragraph> */}
      </Typography>

      <Divider />

      <Row gutter={16}>
        <Col span={8}>
          <Card bordered={false}>
            <Statistic
              title="平台收入金额"
              value={summary.totalAmount}
              valueStyle={{ color: '#3f8600' }}
              suffix="元"
              formatter={formatter}
              loading={isLoading}
            />
          </Card>
        </Col>
        {/* <Col span={8}>
          <Card bordered={false}>
            <Statistic
              title="可用账号"
              value={0}
              valueStyle={{ color: '#3f8600' }}
              suffix="个"
            />
          </Card>
        </Col> */}
        {/* <Col span={8}>
          <Card bordered={false}>
            <Statistic
              title="可用代理IP"
              value={0}
              valueStyle={{ color: '#3f8600' }}
              suffix="个"
            />
          </Card>
        </Col> */}
      </Row>
    </>
  )
}

export default Dashboard