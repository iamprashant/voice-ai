import type { FC, HTMLAttributes } from 'react';
import React, { useEffect, useState } from 'react';
import ReactECharts from 'echarts-for-react';
import type { EChartsOption } from 'echarts';
import dayjs from 'dayjs';
import { get } from 'lodash-es';
import { formatNumber } from '@/utils/format';
import type { AppTokenCostsResponse } from '@/models/app';
import { cn, toDate } from '@/styles/media';
import { Tooltip } from '@/app/components/Tooltip';
import { InfoIcon } from '@/app/components/Icon/Info';
import { Label } from '@/app/components/Form/Label';
import { AuditLog } from '@rapidaai/react';
import { Spinner } from '@/app/components/Loader/Spinner';
import { getTotalTokenMetric } from '@/utils/metadata';
import { Card } from '@/app/components/base/cards';
import { useDarkMode } from '@/context/dark-mode-context';

const valueFormatter = (v: string | number) => v;

const COLOR_TYPE_MAP = {
  green: {
    lineColor: 'rgba(6, 148, 162, 1)',
    bgColor: ['rgba(6, 148, 162, 0.2)', 'rgba(67, 174, 185, 0.08)'],
  },
  orange: {
    lineColor: 'rgba(255, 138, 76, 1)',
    bgColor: ['rgba(254, 145, 87, 0.2)', 'rgba(255, 138, 76, 0.1)'],
  },
  blue: {
    lineColor: 'rgba(28, 100, 242, 1)',
    bgColor: ['rgba(28, 100, 242, 0.3)', 'rgba(28, 100, 242, 0.1)'],
  },
};

const COMMON_COLOR_MAP = {
  label: '#9CA3AF',
  splitLineLight: '#e6e6e6',
  splitLineDark: '#121212',
};

type IColorType = 'green' | 'orange' | 'blue';
type IChartType = 'conversations' | 'endUsers' | 'costs' | 'workflowCosts';
type IChartConfigType = { colorType: IColorType; showTokens?: boolean };

const commonDateFormat = 'MMM D, YYYY';

const CHART_TYPE_CONFIG: Record<string, IChartConfigType> = {
  conversations: {
    colorType: 'green',
  },
  endUsers: {
    colorType: 'orange',
  },
  costs: {
    colorType: 'blue',
    showTokens: true,
  },
  workflowCosts: {
    colorType: 'blue',
  },
};

export type PeriodParams = {
  start: Date;
  end: Date;
} & HTMLAttributes<HTMLDivElement>;

export type IChartProps = {
  valueKey?: string;
  yMax?: number;
  chartType: IChartType;
  height?: number;
  chartData:
    | AppTokenCostsResponse
    | { data: Array<{ date: string; count: number }> };
};

const Chart: React.FC<IChartProps> = ({
  chartType = 'conversations',
  height = 160,
  chartData,
  valueKey,
  yMax,
}) => {
  const { isDarkMode } = useDarkMode();
  const statistics = chartData.data;
  const statisticsLen = statistics.length;
  const extraDataForMarkLine = new Array(
    statisticsLen >= 2 ? statisticsLen - 2 : statisticsLen,
  ).fill('1');
  extraDataForMarkLine.push('');
  extraDataForMarkLine.unshift('');

  const xData = statistics.map(({ date }) => date);
  const yField =
    valueKey ||
    Object.keys(statistics[0]).find(name => name.includes('count')) ||
    '';
  const yData = statistics.map(item => {
    return item[yField] || 0;
  });

  const options: EChartsOption = {
    dataset: {
      dimensions: ['date', yField],
      source: statistics,
    },
    grid: { top: 8, right: 36, bottom: 0, left: 0, containLabel: true },
    tooltip: {
      trigger: 'item',
      position: 'top',
      borderWidth: 0,
    },
    xAxis: [
      {
        type: 'category',
        boundaryGap: false,
        axisLabel: {
          color: COMMON_COLOR_MAP.label,
          hideOverlap: true,
          overflow: 'break',
          formatter(value) {
            return dayjs(value).format(commonDateFormat);
          },
        },
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: {
          show: true,
          lineStyle: {
            color: isDarkMode
              ? COMMON_COLOR_MAP.splitLineDark
              : COMMON_COLOR_MAP.splitLineLight,
            width: 1,
            type: [10, 10],
          },
          interval(index) {
            return index === 0 || index === xData.length - 1;
          },
        },
      },
      {
        position: 'bottom',
        boundaryGap: false,
        axisLabel: { show: false },
        axisLine: { show: false },
        axisTick: { show: false },
        splitLine: {
          show: true,
          lineStyle: {
            color: isDarkMode
              ? COMMON_COLOR_MAP.splitLineDark
              : COMMON_COLOR_MAP.splitLineLight,
          },
          interval(index, value) {
            return !!value;
          },
        },
      },
    ],
    yAxis: {
      max: yMax ?? 'dataMax',
      type: 'value',
      axisLabel: { color: COMMON_COLOR_MAP.label, hideOverlap: true },
      splitLine: {
        lineStyle: {
          color: isDarkMode
            ? COMMON_COLOR_MAP.splitLineDark
            : COMMON_COLOR_MAP.splitLineLight,
        },
      },
    },
    series: [
      {
        type: 'line',
        showSymbol: true,
        smooth: true,
        symbolSize: 4,
        lineStyle: {
          color:
            COLOR_TYPE_MAP[CHART_TYPE_CONFIG[chartType].colorType].lineColor,
          width: 2,
        },
        itemStyle: {
          color:
            COLOR_TYPE_MAP[CHART_TYPE_CONFIG[chartType].colorType].lineColor,
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0,
                color:
                  COLOR_TYPE_MAP[CHART_TYPE_CONFIG[chartType].colorType]
                    .bgColor[0],
              },
              {
                offset: 1,
                color:
                  COLOR_TYPE_MAP[CHART_TYPE_CONFIG[chartType].colorType]
                    .bgColor[1],
              },
            ],
            global: false,
          },
        },
        tooltip: {
          backgroundColor: isDarkMode ? '#272727' : '#efefef',
          borderWidth: 1,
          borderColor: isDarkMode ? '#2e2e2e' : '#b5b5b5',
          padding: [8, 12, 8, 12],
          formatter(params) {
            return `<div class="text-sm">${params.name}</div>
                          <div style='font-size:14px;'>${valueFormatter((params.data as any)[yField])}
                              ${
                                !CHART_TYPE_CONFIG[chartType].showTokens
                                  ? ''
                                  : `<span style='font-size:12px'>
                                  <span style='margin-left:4px;'>(</span>
                                  <span style='color:#FF8A4C'>~$${get(params.data, 'total_price', 0)}</span>
                                  <span style='color:#6B7280'>)</span>
                              </span>`
                              }
                          </div>`;
          },
        },
      },
    ],
  };
  return <ReactECharts option={options} style={{ height: height }} />;
};
export type AppStatisticsResponse = {
  data: Array<{ date: string }>;
};

/**
 *
 * @param props
 * @returns
 */
export const TokenOutputSpeed: FC<
  PeriodParams & {
    loading: boolean;
    activities: AuditLog[];
  }
> = props => {
  useEffect(() => {
    if (props.activities.length > 0) {
      //   const totalTimeTaken = props.activities.reduce(
      //     (accumulator, currentValue) =>
      //       accumulator + currentValue.getMetricsList() / 1000000,
      //     0,
      //   );
    }
  }, [props.activities]);
  return (
    <Card className={props.className}>
      <Label>
        <span className="text-base font-semibold">Token Output Speed</span>
        <Tooltip
          icon={
            <InfoIcon className="w-4 h-4 mt-[2px] ml-0.5dark:text-gray-400" />
          }
        >
          <p className={cn('font-normal text-sm p-1 w-64')}>
            Reflect daily token usage of large language models.
          </p>
        </Tooltip>
      </Label>
      <div className="text-xs font-normalgroup-hover:text-gray-700 break-all ">
        Last{' '}
        {Math.ceil(
          Math.abs(props.end.getTime() - props.start.getTime()) /
            (1000 * 3600 * 24),
        )}{' '}
        Days
      </div>

      <div
        className={cn(
          'my-4 flex-1',
          'flex flex-row items-center break-all text-2xl font-semibold',
          props.activities.length === 0 && 'opacity-60 text-3xl',
        )}
      >
        {props.activities.length > 0
          ? `${formatNumber(87 / 1000)}k Tokens/s`
          : 0}
      </div>
      {props.loading ? (
        <div className="w-full h-[160px] items-center justify-center flex">
          <Spinner size="md" />
        </div>
      ) : (
        <Chart
          height={160}
          chartData={
            {
              data:
                props.activities.length > 0
                  ? getActualChartData(props.start, props.end, props.activities)
                  : getDefaultChartData({
                      start: props.start,
                      end: props.end,
                      key: 'count',
                    }),
            } as any
          }
          chartType="conversations"
          valueKey="count"
          {...{ yMax: 500 }}
        />
      )}
    </Card>
  );
};

/**
 *
 */
export const AverageResponseTime: FC<
  PeriodParams & {
    loading: boolean;
    activities: AuditLog[];
  }
> = props => {
  //   const response: AppStatisticsResponse = { data: [] };

  const [latency, setLatency] = useState(0);
  const [maxLatency, setMaxLatency] = useState(900);
  useEffect(() => {
    if (props.activities.length > 0) {
      const totalTimeTaken = props.activities.reduce(
        (accumulator, currentValue) =>
          accumulator + currentValue.getTimetaken() / 1000000,
        0,
      );
      setLatency(totalTimeTaken / props.activities.length);
      setMaxLatency(20000);
    }
  }, [props.activities]);

  return (
    <Card className={props.className}>
      <Label>
        <span className="text-base font-semibold">Avg. Response Time</span>
        <Tooltip
          icon={
            <InfoIcon className="w-4 h-4 mt-[2px] ml-0.5dark:text-gray-400" />
          }
        >
          <p className={cn('font-normal text-sm p-1 w-64')}>
            Reflect daily token usage of large language models.
          </p>
        </Tooltip>
      </Label>
      <div className="text-xs font-normalgroup-hover:text-gray-700 break-all ">
        Last{' '}
        {Math.ceil(
          Math.abs(props.end.getTime() - props.start.getTime()) /
            (1000 * 3600 * 24),
        )}{' '}
        Days
      </div>
      <div
        className={cn(
          'my-4 flex-1',
          'flex flex-row items-center break-all text-2xl font-semibold',
          props.activities.length === 0 && 'opacity-30 text-3xl',
        )}
      >
        {props.activities.length > 0
          ? `~ ${formatNumber(Math.ceil(latency))}ms/request`
          : 0}
      </div>
      {props.loading ? (
        <div className="w-full h-[300px] items-center justify-center flex">
          <Spinner size="md" />
        </div>
      ) : (
        <Chart
          chartData={
            {
              data:
                props.activities.length > 0
                  ? getActualChartData(props.start, props.end, props.activities)
                  : getDefaultChartData({
                      start: props.start,
                      end: props.end,
                      key: 'count',
                    }),
            } as any
          }
          height={300}
          chartType="conversations"
          valueKey="count"
          {...{ yMax: maxLatency }}
        />
      )}
    </Card>
  );
};

/**
 *
 * @param period
 * @returns
 */

export const TokenUsages: FC<
  PeriodParams & {
    loading: boolean;
    activities: AuditLog[];
  }
> = ({ start, end, className, activities, loading }) => {
  const response: AppStatisticsResponse = { data: [] };
  const noDataFlag = !response.data || response.data.length === 0;

  const [totalToken, setTotalToken] = useState(0);
  const [maxToken, setMaxToken] = useState(500);
  useEffect(() => {
    if (activities.length > 0) {
      let tokens = 0;
      let _maxToken = 500;
      activities.forEach(x => {
        let currentToken = getTotalTokenMetric(x.getMetricsList());
        tokens += currentToken;
        _maxToken = _maxToken < currentToken ? currentToken : _maxToken;
      });
      setTotalToken(tokens);
      setMaxToken(_maxToken + 100);
    }
  }, [activities]);

  return (
    <Card className={className}>
      <Label>
        <span className="text-base font-semibold">Token Usage</span>
        <Tooltip
          icon={
            <InfoIcon className="w-4 h-4 mt-[2px] ml-0.5dark:text-gray-400" />
          }
        >
          <p className={cn('font-normal text-sm p-1 w-64')}>
            Reflect daily token usage of large language models.
          </p>
        </Tooltip>
      </Label>
      <div className="text-xs font-normalgroup-hover:text-gray-700 break-all ">
        Last{' '}
        {Math.ceil(
          Math.abs(end.getTime() - start.getTime()) / (1000 * 3600 * 24),
        )}{' '}
        Days
      </div>
      <div
        className={cn(
          'my-4 flex-1',
          'flex flex-row items-center break-all text-2xl font-semibold',
          activities.length === 0 && 'opacity-60 text-3xl',
        )}
      >
        {activities.length > 0
          ? `${formatNumber(totalToken / 1000)}k Tokens`
          : 0}
      </div>
      {loading ? (
        <div className="w-full h-[160px] items-center justify-center flex">
          <Spinner size="md" />
        </div>
      ) : (
        <Chart
          height={160}
          chartData={
            {
              data:
                activities.length > 0
                  ? getActualChartData(start, end, activities)
                  : getDefaultChartData({
                      start: start,
                      end: end,
                      key: 'tokens',
                    }),
            } as any
          }
          chartType="conversations"
          valueKey="tokens"
          {...(noDataFlag && { yMax: maxToken })}
        />
      )}
    </Card>
  );
};

const getDefaultChartData = ({
  start,
  end,
  key = 'count',
}: {
  start: Date;
  end: Date;
  key?: string;
}) => {
  //   const diffDays = Math.max(dayjs(end).diff(dayjs(start), 'day'), 1);
  const diffDays = Math.ceil(
    Math.abs(end.getTime() - start.getTime()) / (1000 * 3600 * 24),
  );
  return Array.from({ length: diffDays }, (_, index) => ({
    date: dayjs(end).add(index, 'day').format(commonDateFormat),
    [key]: 0,
  }));
};

/**
 *
 * @param models
 * @returns
 */
function prepareChartData(models: AuditLog[]) {
  const dataMap = {}; // Map to store aggregated timetaken for each date

  // Iterate over each model
  models.forEach(model => {
    // Extract relevant data from the model
    let date = model.getCreateddate();
    // If date is not present in the model, set it to 0
    if (!date) return;

    // Parse date using dayjs
    const formattedDate = `${dayjs(toDate(date)).format(commonDateFormat)}`;

    // If date is not present in dataMap, initialize it
    if (!dataMap[formattedDate]) {
      dataMap[formattedDate] = {
        totalTimetaken: 0,
        totalToken: 0,
        count: 0,
      };
    }

    // Accumulate timetaken for the date
    dataMap[formattedDate].totalTimetaken += model.getTimetaken();
    dataMap[formattedDate].totalToken += getTotalTokenMetric(
      model.getMetricsList(),
    );
    dataMap[formattedDate].count++;
  });
  return dataMap;
}

/**
 *
 * @param start
 * @param end
 * @param models
 * @returns
 */
function getActualChartData(start: Date, end: Date, models: AuditLog[]) {
  // Calculate average timetaken for each date and construct the output array

  const dataMap = prepareChartData(models);
  const diffDays = Math.ceil(
    Math.abs(end.getTime() - start.getTime()) / (1000 * 3600 * 24),
  );
  return Array.from({ length: diffDays }, (_, index) => {
    const currentDate = dayjs(start).add(index, 'day').format(commonDateFormat);
    const aggregatedData = dataMap[currentDate];
    // Calculate average timetaken or set to 0 if no data available
    const averageTimetaken = aggregatedData
      ? Math.ceil(
          aggregatedData.totalTimetaken / (aggregatedData.count * 1000000),
        )
      : 0;

    return {
      date: currentDate,
      count: averageTimetaken,
      tokens: aggregatedData ? aggregatedData.totalToken : 0,
    };
  });
}

export default Chart;
