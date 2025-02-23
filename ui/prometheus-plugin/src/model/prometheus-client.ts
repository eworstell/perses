// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { fetch, fetchJson, RequestHeaders } from '@perses-dev/core';
import { DatasourceClient } from '@perses-dev/plugin-system';
import {
  InstantQueryRequestParameters,
  InstantQueryResponse,
  LabelNamesRequestParameters,
  LabelNamesResponse,
  LabelValuesRequestParameters,
  LabelValuesResponse,
  RangeQueryRequestParameters,
  RangeQueryResponse,
} from './api-types';

interface PrometheusClientOptions {
  datasourceUrl: string;
  headers?: RequestHeaders;
}

export interface PrometheusClient extends DatasourceClient {
  options: PrometheusClientOptions;
  instantQuery(params: InstantQueryRequestParameters, headers?: RequestHeaders): Promise<InstantQueryResponse>;
  rangeQuery(params: RangeQueryRequestParameters, headers?: RequestHeaders): Promise<RangeQueryResponse>;
  labelNames(params: LabelNamesRequestParameters, headers?: RequestHeaders): Promise<LabelNamesResponse>;
  labelValues(params: LabelValuesRequestParameters, headers?: RequestHeaders): Promise<LabelValuesResponse>;
}

export interface QueryOptions {
  datasourceUrl: string;
  headers?: RequestHeaders;
}

/**
 * Calls the `/api/v1/query` endpoint to get metrics data.
 */
export function instantQuery(params: InstantQueryRequestParameters, queryOptions: QueryOptions) {
  return fetchWithPost<InstantQueryRequestParameters, InstantQueryResponse>('/api/v1/query', params, queryOptions);
}

/**
 * Calls the `/api/v1/query_range` endpoint to get metrics data.
 */
export function rangeQuery(params: RangeQueryRequestParameters, queryOptions: QueryOptions) {
  return fetchWithPost<RangeQueryRequestParameters, RangeQueryResponse>('/api/v1/query_range', params, queryOptions);
}

/**
 * Calls the `/api/v1/labels` endpoint to get a list of label names.
 */
export function labelNames(params: LabelNamesRequestParameters, queryOptions: QueryOptions) {
  return fetchWithPost<LabelNamesRequestParameters, LabelNamesResponse>('/api/v1/labels', params, queryOptions);
}

/**
 * Calls the `/api/v1/label/{labelName}/values` endpoint to get a list of values for a label.
 */
export function labelValues(params: LabelValuesRequestParameters, queryOptions: QueryOptions) {
  const { labelName, ...searchParams } = params;
  const apiURI = `/api/v1/label/${encodeURIComponent(labelName)}/values`;
  return fetchWithGet<typeof searchParams, LabelValuesResponse>(apiURI, searchParams, queryOptions);
}

function fetchWithGet<T extends RequestParams<T>, TResponse>(apiURI: string, params: T, queryOptions: QueryOptions) {
  const { datasourceUrl } = queryOptions;

  let url = `${datasourceUrl}${apiURI}`;
  const urlParams = createSearchParams(params).toString();
  if (urlParams !== '') {
    url += `?${urlParams}`;
  }
  return fetchJson<TResponse>(url, { method: 'GET' });
}

function fetchWithPost<T extends RequestParams<T>, TResponse>(apiURI: string, params: T, queryOptions: QueryOptions) {
  const { datasourceUrl, headers } = queryOptions;

  const url = `${datasourceUrl}${apiURI}`;
  const init = {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      ...headers,
    },
    body: createSearchParams(params),
  };
  return fetchResults<TResponse>(url, init);
}

// Request parameter values we know how to serialize
type ParamValue = string | string[] | number | undefined;

// Used to constrain the types that can be passed to createSearchParams to
// just the ones we know how to serialize
type RequestParams<T> = {
  [K in keyof T]: ParamValue;
};

/**
 * Creates URLSearchParams from a request params object.
 */
function createSearchParams<T extends RequestParams<T>>(params: T) {
  const searchParams = new URLSearchParams();
  for (const key in params) {
    const value: ParamValue = params[key];
    if (value === undefined) continue;

    if (typeof value === 'string') {
      searchParams.append(key, value);
      continue;
    }

    if (typeof value === 'number') {
      searchParams.append(key, value.toString());
      continue;
    }

    for (const val of value) {
      searchParams.append(key, val);
    }
  }
  return searchParams;
}

/**
 * Fetch JSON and parse warnings for query inspector
 */
export async function fetchResults<T>(...args: Parameters<typeof global.fetch>) {
  const response = await fetch(...args);
  const json: T = await response.json();
  return { ...json, rawResponse: response };
}
