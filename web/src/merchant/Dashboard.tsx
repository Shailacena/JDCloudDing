import { Card, Col, Row, Statistic, StatisticProps, Typography } from "antd";
import { useEffect, useMemo, useState } from "react";
import { getDataFormat, greet } from "../utils/Tool";
import CountUp from 'react-countup';
import { Line } from "@ant-design/charts";
import { BaseDailyBill, ListDailyBillReq } from "../api/types";
import { useApis } from "../api/api";
import { useAppContext } from "../AppProvider";

interface DataType extends BaseDailyBill {
  '日期': string;
  '收入': number;
}

const Dashboard: React.FC = () => {
  const [data, setData] = useState<DataType[]>([])
  const [todayAmount, setTodayAmount] = useState(0)
  const [todayNum, setTodayNum] = useState(0)
  const [twoWeekAmount, setTwoWeekAmount] = useState(0)
  const [twoWeekNum, setTwoWeekNum] = useState(0)
  const [balance, setBalance] = useState(0)
  let { listMerchant1StatisticsBill, merchantGetBalance } = useApis()
  const ctx = useAppContext();
  const [currentTime, setCurrentTime] = useState(getDataFormat(new Date()));

  const fetchListStatisticsBill = async () => {

    let params: ListDailyBillReq = {
      merchantId: ctx.cookie.id,
    }
    const { data } = await listMerchant1StatisticsBill(params)
    let totalAmount = 0
    let totalNum = 0
    const date = new Date();
    const todayStr = date.toISOString().slice(0, 10);
    let d: DataType[] = data?.list?.map((item, _) => {
      let newItem: DataType = {
        '收入': item.totalSuccessAmount,
        '日期': item.date,
        ...item
      }
      totalAmount += item.totalSuccessAmount;
      totalNum += item.totalOrderNum
      if (todayStr == item.date) {
        setTodayAmount(item.totalSuccessAmount)
        setTodayNum(item.totalOrderNum)
      }
      return newItem
    })
    setTwoWeekAmount(totalAmount)
    setTwoWeekNum(totalNum)
    d.reverse()
    setData(d)
  }

  useEffect(() => {
    fetchListStatisticsBill()
    fetchListMerchant()
  }, [])

  const fetchListMerchant = async () => {
    const { data } = await merchantGetBalance()
    setBalance(data.balance)
  }

  const formatter: StatisticProps['formatter'] = (value) => (
    <CountUp end={value as number} separator="," />
  );

  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentTime(getDataFormat(new Date()));
    }, 1000);
    // 清除间隔，以防内存泄漏
    return () => clearInterval(interval);
  }, []);

  const CurrentTime = () => {
    return (
      <Typography.Text>{currentTime}</Typography.Text>
    );
  };

  const ChartLine: React.FC = () => {
    const config = {
      data,
      smooth: true,
      height: 220,
      xField: '日期',
      yField: '收入',
      seriesField: 'category',
      point: {
        shapeField: 'circle',
        sizeField: 4,
      },
      interaction: {
        tooltip: {
          marker: true,
        },
      },
      style: {
        lineWidth: 2,
      },
    };
    return <Line {...config} />;
  };

  return (
    <div style={{ margin: '0 auto' }}>
      {/* <CustomBreadcrumb arr={['基本','按钮']}/> */}
      <div className="div_time">
        <p>{greet()}, 老板! &nbsp;&nbsp;现在时间是: {CurrentTime()}</p>
      </div>
      <Row gutter={16}>
        <Col span={4}>
          <Card title="余额">
            <Row gutter={16}>
              <Col span={24}>
                <Statistic title="--" value={balance} formatter={formatter} />
              </Col>
            </Row>
          </Card>
        </Col>
        <Col span={10}>
          <Card title="今日">
            <Row gutter={16}>
              <Col span={12}>
                <Statistic title="营收" value={todayAmount} formatter={formatter} />
              </Col>
              <Col span={12}>
                <Statistic title="订单" value={todayNum} precision={2} formatter={formatter} />
              </Col>
            </Row>
          </Card>
        </Col>
        <Col span={10}>
          <Card title="最近两周">
            <Row gutter={16}>
              <Col span={12}>
                <Statistic title="营收" value={twoWeekAmount} formatter={formatter} />
              </Col>
              <Col span={12}>
                <Statistic title="订单" value={twoWeekNum} precision={2} formatter={formatter} />
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>

      <br></br>

      <Card title="最近14天">
        {
          useMemo(() => <ChartLine />, [data])
        }

      </Card>
      <div>
      </div>
    </div>
  );
};

export default Dashboard

