import { useEffect, useMemo, useState } from 'react';
import { Card, Col, Row, Statistic, Typography, Space, Button, Switch, Segmented, Select } from 'antd';
import { Line } from '@ant-design/charts';
import { useApis } from '../api/api';
import { GetServerStatusResp, ListPartnerReq, ListMerchantReq, TrendPoint } from '../api/types';

const { Title, Text } = Typography;

function ServerStatus() {
  const { adminGetServerStatus, adminGetTodayTrend, listPartner, listMerchant } = useApis()

  const [status, setStatus] = useState<GetServerStatusResp>()
  

  // 数据类型：金额 / 数量
  const [metric, setMetric] = useState<'amount' | 'count'>('count')

  // 主账号筛选：合作商 / 商户
  const [accountType, setAccountType] = useState<'partner' | 'merchant' | 'all'>('all')
  const [partnerOptions, setPartnerOptions] = useState<{ value: number, label: string }[]>([])
  const [merchantOptions, setMerchantOptions] = useState<{ value: number, label: string }[]>([])
  const [selectedPartnerIds, setSelectedPartnerIds] = useState<number[]>([])
  const [selectedMerchantIds, setSelectedMerchantIds] = useState<number[]>([])

  // 对比数据与粒度
  const [seriesData, setSeriesData] = useState<Array<{ label: string; points: TrendPoint[] }>>([])
  const [granularity, setGranularity] = useState<5 | 10 | 15>(5)

  // 自动刷新
  const [autoRefresh, setAutoRefresh] = useState<boolean>(false)
  const [refreshInterval, setRefreshInterval] = useState<number>(60) // 秒

  useEffect(() => {
    refresh()
    // 预加载主账号列表
    fetchPartners()
    fetchMerchants()
  }, [])

  const refresh = async () => {
    const res1 = await adminGetServerStatus({})
    if (res1?.success) setStatus(res1.data)

    // 组装多账号趋势数据
    const series: Array<{ label: string; points: TrendPoint[] }> = []

    const buildLabel = (type: 'partner' | 'merchant', id: number) => {
      const opts = type === 'partner' ? partnerOptions : merchantOptions
      const found = opts.find(o => o.value === id)
      const name = found?.label || String(id)
      return type === 'partner' ? `合作商-${name}` : `商户-${name}`
    }

    if (accountType === 'all') {
      const res2 = await adminGetTodayTrend({})
      if (res2?.success) {
        series.push({ label: '全部', points: res2.data.points || [] })
      }
    } else if (accountType === 'partner') {
      const ids = selectedPartnerIds?.length ? selectedPartnerIds : []
      if (ids.length === 0) {
        const res2 = await adminGetTodayTrend({})
        if (res2?.success) {
          series.push({ label: '全部', points: res2.data.points || [] })
        }
      } else {
        const results = await Promise.all(ids.map(id => adminGetTodayTrend({ partnerId: id })))
        results.forEach((r, idx) => {
          if (r?.success) series.push({ label: buildLabel('partner', ids[idx]), points: r.data?.points || [] })
        })
      }
    } else if (accountType === 'merchant') {
      const ids = selectedMerchantIds?.length ? selectedMerchantIds : []
      if (ids.length === 0) {
        const res2 = await adminGetTodayTrend({})
        if (res2?.success) {
          series.push({ label: '全部', points: res2.data.points || [] })
        }
      } else {
        const results = await Promise.all(ids.map(id => adminGetTodayTrend({ merchantId: id })))
        results.forEach((r, idx) => {
          if (r?.success) series.push({ label: buildLabel('merchant', ids[idx]), points: r.data?.points || [] })
        })
      }
    }

    setSeriesData(series)
  }

  const fetchPartners = async () => {
    try {
      const { data } = await listPartner({ ignoreStatistics: true } as ListPartnerReq)
      const opts = (data?.list || []).map((p: any) => ({ value: p.id, label: `${p.id}(${p.nickname})` }))
      setPartnerOptions(opts)
    } catch (e) { /* ignore */ }
  }

  const fetchMerchants = async () => {
    try {
      const { data } = await listMerchant({ ignoreStatistics: true } as ListMerchantReq)
      const opts = (data?.list || []).map((m: any) => ({ value: m.id, label: `${m.id}(${m.nickname || m.username})` }))
      setMerchantOptions(opts)
    } catch (e) { /* ignore */ }
  }

  const toMinutes = (hhmm: string): number => {
    const [h, m] = hhmm.split(':')
    const hh = parseInt(h || '0', 10)
    const mm = parseInt(m || '0', 10)
    return hh * 60 + mm
  }

  const aggregatePoints = (points: TrendPoint[], g: number): TrendPoint[] => {
    if (g === 5) return points
    const bucketMap = new Map<number, TrendPoint>()
    const bucketSize = g
    points.forEach(p => {
      const minutes = toMinutes(p.time)
      const bucket = Math.floor(minutes / bucketSize) * bucketSize
      const key = bucket
      const existing = bucketMap.get(key)
      if (!existing) {
        bucketMap.set(key, {
          time: `${String(Math.floor(bucket / 60)).padStart(2, '0')}:${String(bucket % 60).padStart(2, '0')}`,
          orderCount: p.orderCount,
          orderAmount: p.orderAmount,
          successCount: p.successCount,
          successAmount: p.successAmount,
        })
      } else {
        existing.orderCount += p.orderCount
        existing.orderAmount += p.orderAmount
        existing.successCount += p.successCount
        existing.successAmount += p.successAmount
      }
    })
    return Array.from(bucketMap.entries()).sort((a, b) => a[0] - b[0]).map(([, v]) => v)
  }

  const chartData = useMemo(() => {
    const d: any[] = []
    seriesData.forEach(series => {
      const agg = aggregatePoints(series.points, granularity)
      agg.forEach(p => {
        if (metric === 'amount') {
          d.push({ 日期: p.time, 数值: p.orderAmount, category: `订单金额-${series.label}` })
          d.push({ 日期: p.time, 数值: p.successAmount, category: `到账金额-${series.label}` })
        } else {
          d.push({ 日期: p.time, 数值: p.orderCount, category: `订单数-${series.label}` })
          d.push({ 日期: p.time, 数值: p.successCount, category: `成功数-${series.label}` })
        }
      })
    })
    return d
  }, [seriesData, metric, granularity])

  const chartConfig = {
    data: chartData,
    smooth: true,
    height: 260,
    xField: '日期',
    yField: '数值',
    seriesField: 'category',
    point: { shapeField: 'circle', sizeField: 3 },
    interaction: { tooltip: { marker: true } },
    style: { lineWidth: 2 },
    slider: { start: 0.5, end: 1 },
  }

  const exportCSV = () => {
    const rows: string[] = []
    rows.push(['时间', '账号', '订单数', '成功数', '订单金额', '到账金额'].join(','))
    seriesData.forEach(series => {
      const agg = aggregatePoints(series.points, granularity)
      agg.forEach(p => {
        rows.push([
          p.time,
          series.label,
          String(p.orderCount),
          String(p.successCount),
          String(p.orderAmount),
          String(p.successAmount),
        ].join(','))
      })
    })
    const blob = new Blob([rows.join('\n')], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `today_trend_${metric}_${granularity}min.csv`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }

  useEffect(() => {
    if (!autoRefresh) return
    const timer = setInterval(() => { refresh() }, refreshInterval * 1000)
    return () => clearInterval(timer)
  }, [autoRefresh, refreshInterval])

  return (
    <>
      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card>
            <Title level={4}>数据库状态</Title>
            <Row gutter={[16, 16]}>
              <Col span={4}><Statistic title="数据库版本" value={status?.dbVersion || '-'} /></Col>
              <Col span={4}><Statistic title="订单总量" value={status?.orderTotal || 0} /></Col>
              <Col span={4}><Statistic title="今日订单数" value={status?.orderToday || 0} /></Col>
              <Col span={4}><Statistic title="今日成功数" value={status?.successTodayCount || 0} /></Col>
              <Col span={4}><Statistic title="今日到账金额" precision={2} value={status?.successTodayAmount || 0} /></Col>
              <Col span={2}><Statistic title="合作商" value={status?.partnerTotal || 0} /></Col>
              <Col span={2}><Statistic title="商户" value={status?.merchantTotal || 0} /></Col>
            </Row>
          </Card>
        </Col>

        <Col span={24}>
          <Card>
            <Space direction="vertical" style={{ width: '100%' }}>
              <Title level={4}>当日订单趋势（5分钟）</Title>
              <Space wrap>
                <Text>数据类型</Text>
                <Segmented options={[{ label: '数量', value: 'count' }, { label: '金额', value: 'amount' }]} value={metric} onChange={(v) => setMetric(v as any)} />
                <Text>粒度</Text>
                <Segmented options={[{ label: '5分钟', value: 5 }, { label: '10分钟', value: 10 }, { label: '15分钟', value: 15 }]} value={granularity} onChange={(v) => setGranularity(v as any)} />
                <Text>主账号类型</Text>
                <Segmented options={[{ label: '全部', value: 'all' }, { label: '合作商', value: 'partner' }, { label: '商户', value: 'merchant' }]} value={accountType} onChange={(v) => {
                  setAccountType(v as any)
                  setSelectedPartnerIds([])
                  setSelectedMerchantIds([])
                  setTimeout(() => { refresh() }, 0)
                }} />
                {accountType === 'partner' && (
                  <>
                    <Text>合作商</Text>
                    <Select mode="multiple" allowClear maxTagCount={3} style={{ minWidth: 320 }} options={partnerOptions} value={selectedPartnerIds} onChange={(vals) => { setSelectedPartnerIds(vals as number[]); setTimeout(() => { refresh() }, 0) }} />
                  </>
                )}
                {accountType === 'merchant' && (
                  <>
                    <Text>商户</Text>
                    <Select mode="multiple" allowClear maxTagCount={3} style={{ minWidth: 320 }} options={merchantOptions} value={selectedMerchantIds} onChange={(vals) => { setSelectedMerchantIds(vals as number[]); setTimeout(() => { refresh() }, 0) }} />
                  </>
                )}
                <Text>自动刷新</Text>
                <Switch checked={autoRefresh} onChange={setAutoRefresh} />
                <Text>间隔</Text>
                <Select style={{ width: 120 }} value={refreshInterval} onChange={(v) => setRefreshInterval(v)} options={[{ value: 30, label: '30秒' }, { value: 60, label: '60秒' }, { value: 300, label: '5分钟' }]} />
                <Button onClick={refresh}>手动刷新</Button>
                <Button onClick={exportCSV}>导出 CSV</Button>
              </Space>
              <Line {...chartConfig} />
            </Space>
          </Card>
        </Col>
      </Row>
    </>
  )
}

export default ServerStatus