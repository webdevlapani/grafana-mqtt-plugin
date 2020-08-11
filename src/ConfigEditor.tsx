import _ from 'lodash';
import * as $ from 'jquery';
import React, { PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import { DataSourceOptions, Regions } from './types';

const { FormField, Select } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<DataSourceOptions> {}

interface State {
  region: SelectableValue<string>;
  endpoint: string;
  certificates: any[];
}

export class ConfigEditor extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);
    const region = _.find(Regions, { value: props.options.jsonData.region }) || Regions[0];
    this.state = {
      region: region,
      endpoint: '',
      certificates: [],
    };

    this.reload(region);
  }

  onRegionChange = (value: SelectableValue<string>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
    };
    jsonData.region = value.value as string;
    onOptionsChange({ ...options, jsonData });

    this.setState({
      region: value,
      endpoint: '',
      certificates: [],
    });

    this.reload(value);
  };

  reload = (region: SelectableValue<string>) => {
    const { options } = this.props;

    $.get(`/api/datasources/${options.id}/resources/endpoint?region=${region.value}`)
      .then(data => {
        this.setState({ endpoint: data });
      })
      .catch(err => {
        // TODO: alert
        console.log(err);
      });

    $.get(`/api/datasources/${options.id}/resources/certificates?region=${region.value}`)
      .then(data => {
        this.setState({ certificates: JSON.parse(data) });
      })
      .catch(err => {
        // TODO: alert
        console.log(err);
      });
  };

  render() {
    const { region, endpoint, certificates } = this.state;

    return [
      <div className="gf-form-group" key="region">
        <div className="gf-form">
          <FormField
            label="Region"
            labelWidth={10}
            required
            inputEl={<Select width={27} options={Regions} value={region} onChange={this.onRegionChange} />}
          />
        </div>
      </div>,
      <div className="gf-form-group" key="configuration">
        <div className="gf-form">
          <FormField
            label="MQTT Host"
            labelWidth={10}
            inputWidth={27}
            value={endpoint}
            placeholder="MQTT Host"
            readOnly
          />
        </div>
        <div className="gf-form">
          <FormField label="MQTT Port" labelWidth={10} inputWidth={27} value="8883" readOnly />
        </div>
      </div>,
      <div className="gf-form-group" key="certificates">
        {certificates.map(certificate => (
          <p key={certificate.id}>{JSON.stringify(certificate)}</p>
        ))}
      </div>,
    ];
  }
}
