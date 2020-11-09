import _ from 'lodash';
import * as $ from 'jquery';
import React, { PureComponent } from 'react';
import {LegacyForms, Button, Modal} from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, SelectableValue } from '@grafana/data';
import {DataSourceOptions, Regions} from './types';
import axios from 'axios';
import './ConfigEditor.scss';
import CertificateList from "./CertificateList";

const { FormField, Select } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<DataSourceOptions> {}

interface State {
  region: SelectableValue<string>;
  endpoint: string;
  certificates: any[];
  isCertificateModalOpen: boolean;
  certificateModalTitle: string;
  modalBody: string;
  inputTopicPrefix: string;
  inputClientIdPrefix: string;
  newCertificateData: any,
}

export class ConfigEditor extends PureComponent<Props, State> {
  constructor(props: Props) {
    super(props);
    const region = _.find(Regions, { value: props.options.jsonData.region }) || Regions[0];
    this.state = {
      region: region,
      endpoint: '',
      certificates: [],
      isCertificateModalOpen: false,
      certificateModalTitle: 'Create a certificate',
      modalBody: 'createCertificateSection',
      inputTopicPrefix: '*',
      inputClientIdPrefix: '*',
      newCertificateData:{}
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
      isCertificateModalOpen: false,
      certificateModalTitle: 'Create a certificate',
      modalBody: 'createCertificateSection',
      inputTopicPrefix: '*',
      inputClientIdPrefix: '*',
      newCertificateData:{}
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

  /**
   * certificate action
   * @param certificateId
   * @param actionName
   */
  onClickCertificateAction = (certificateId: number, actionName: string) => {
    const { options } = this.props;
    if(actionName === 'disable') {
      this.certificateActionApiCalls(options.id,certificateId, 'set-inactive');
    } else if( actionName === 'enable') {
      this.certificateActionApiCalls(options.id,certificateId, 'set-active')
    } else if( actionName === 'revoked') {
      this.certificateActionApiCalls(options.id,certificateId, 'revoke')
    } else if( actionName === 'delete') {
      this.certificateActionApiCalls(options.id,certificateId, 'delete')
    }
  }

  /**
   * certificate Action Api Calls
   * @param datasourceId
   * @param certificateId
   * @param actionName
   */
  certificateActionApiCalls = (datasourceId: number, certificateId: number, actionName: string ) => {
    axios({
      method: actionName === 'delete' ? 'delete' : 'patch',
      url: `/api/datasources/${datasourceId}/resources/certificates/${actionName}?id=${certificateId}&region=${this.state.region.value}`,

    }).then((data:any)  => {
      this.setState({
        newCertificateData: data
      })
    }).catch( error => {
        // TODO: alert
        console.log(error)
    });
  }

  onClickCreateCertificate = () => {
    this.setState({
      certificateModalTitle: 'Create a policy',
      modalBody: 'createPolicySection'
    })
  }

  onAddCertificateClick = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.preventDefault();
    this.setState({isCertificateModalOpen:true});
  }

  onClickCertificateCreate = () => {
    const {inputClientIdPrefix, inputTopicPrefix} = this.state;
    const { options } = this.props;

    axios({
      method: 'post',
      url: `api/datasources/${options.id}/resources/certificates/create?region=${this.state.region.value}`,
      data: {
        "topic" : inputTopicPrefix,
        "client" : inputClientIdPrefix
      }
    }).then((certificateData: any)=> {
      this.setState({
        newCertificateData:certificateData.data,
        inputClientIdPrefix: '*',
        inputTopicPrefix: '*',
        certificateModalTitle: 'Certificate Created!',
        modalBody:'certificateCreated'
      })
    }).catch((error => {
      // TODO: alert
      console.log(error);
    }));
  }

  createCertificateSection = () => <>
       <p>
         A certificate is used to authenticate your device's connection to AWS IoT.
       </p>
       <div className="types-of-certificate">
         <h2>One-click certificate creation</h2>
         <p>This will generate a certificate, public key, and private key.</p>
         <Button size="md" variant="destructive" onClick={this.onClickCreateCertificate}>
           Create Certificate
         </Button>
       </div>

      <div className="types-of-certificate disable-certificate">
        <h2>Use my certificate <small><b>Coming soon</b></small></h2>
        <p>Use your own certificates for one or many devices.</p>
        <Button size="md" variant="destructive">
          Upload certificate
        </Button>
      </div>

      <div className="types-of-certificate disable-certificate">
        <h2>Register my CA <small><b>Coming soon</b></small></h2>
        <p>Register your CA certificate and enable auto registration of certificates signed by your CA</p>
        <Button size="md" variant="destructive">
          Get Started
        </Button>
      </div>
     </>

  onKeyDownload = (fileType: string) => {
       const { id, public_key, private_key, certificate, root_ca } = this.state.newCertificateData;
       let fileName='';
       let data='';
       if(fileType === 'certificate') {
         fileName = `${id}.cert.pem`
         data=certificate
       } else if(fileType === 'publickey') {
         fileName = `${id}.public.pem`
         data=public_key
       } else if(fileType === 'privatekey') {
         fileName = `${id}.private.pem`
         data=private_key
       } else if(fileType === 'rootca') {
         fileName = `RootCA.pem`
         data=root_ca
       }
    let link = document.createElement('a');
    link.download = fileName;
    let blob = new Blob([data], {type: 'text/plain'});
    link.href = window.URL.createObjectURL(blob);
    link.click();
  }

  viewNewCertificateSection = () => {
     const { id } = this.state.newCertificateData;
       return <>
         <p>
           Download these files and save them in a safe place. These file cannot  be retrieved after you close this page.
         </p>
         <div>
           <b>In order to connect a device, you need to download the following:</b>
           <div className="certificate-download-section">
             <p> A Certificate: </p>
             <p> {`${id}.cert.pem`} </p>
             <Button size="md" variant="destructive" onClick={() =>this.onKeyDownload('certificate')}>
               Download
             </Button>
           </div>
           <div className="certificate-download-section">
             <p> A Public Key: </p>
             <p> {`${id}.public.key`} </p>
             <Button size="md" variant="destructive" onClick={() =>this.onKeyDownload('publickey')}>
               Download
             </Button>
           </div>
           <div className="certificate-download-section">
             <p> A Private Key: </p>
             <p> {`${id}.private.key`} </p>
             <Button size="md" variant="destructive" onClick={() =>this.onKeyDownload('privatekey')}>
               Download
             </Button>
           </div>
           <div className="certificate-download-section">
             <p> A root CA: </p>
             <p> {`RootCA.pem`} </p>
             <Button size="md" variant="destructive" onClick={() =>this.onKeyDownload('rootca')}>
               Download
             </Button>
           </div>
         </div>
         <div>
           <Button className="button-done-certificate" size="md" variant="primary" onClick={() =>this.setState({'isCertificateModalOpen': false})}>
             Done
           </Button>
         </div>
       </>
  }


  createPolicySection = () => <>
    <p>Create a policy to define a set of MQTT topics and client ids. Device with this certificate will only be able to these topic and client ids.</p>
    <div className="gf-form-group" key="configuration">
      <div className="gf-form">
        <FormField
          label="Topic Prefix: 1/2/"
          labelWidth={12}
          inputWidth={30}
          value={this.state.inputTopicPrefix}
          onChange={ (e) => this.setState({inputTopicPrefix: e.target.value})}
        />
      </div>
      <div className="gf-form">
        <FormField
          label="Client ID Prefix: 1/2/"
          labelWidth={12}
          inputWidth={30}
          value={this.state.inputClientIdPrefix}
          onChange={ (e) => this.setState({inputClientIdPrefix: e.target.value})}
        />
      </div>
        <Button
          className="create-certificate-add"
          size="md"
          variant="destructive"
          onClick={this.onClickCertificateCreate}>
          Create
        </Button>
    </div>
  </>

  createCertificateModal = () => {
    const {
      certificateModalTitle,
      modalBody
    } = this.state;

    return <Modal
      isOpen={true}
      onDismiss={() => this.setState({
        isCertificateModalOpen: false,
        modalBody: 'createCertificateSection',
        certificateModalTitle: 'Create a certificate',
      })}
      title={certificateModalTitle}>
      { modalBody === 'createCertificateSection' ?
        this.createCertificateSection() :
        modalBody === 'createPolicySection'?
          this.createPolicySection() :
        modalBody === 'certificateCreated' ?
          this.viewNewCertificateSection(): ''
      }
    </Modal>
  }

  render() {
    const {
      region,
      endpoint,
      certificates,
      isCertificateModalOpen
    } = this.state;

    return <div className="config-editor">
      <div className="gf-form-group" key="region">
        <div className="gf-form">
          <FormField
            label="Region"
            labelWidth={10}
            required
            inputEl={<Select width={27} options={Regions} value={region} onChange={this.onRegionChange} />}
          />
        </div>
      </div>
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
          <FormField label="MQTT Port123" labelWidth={10} inputWidth={27} value="8883" readOnly />
        </div>
      </div>
      <div className="gf-form-group">
        <Button
          icon="plus"
          size="md"
          className="btn-add-certificate"
          variant="destructive"
          onClick={(e) => this.onAddCertificateClick(e)}>
          Add Certificate
        </Button>
      </div>

      <div className="gf-form-group" key="certificates">
        {certificates.map(certificate => (
          <>
            <CertificateList certificate={certificate}>
              {
                certificate.status !== 'REVOKED' ?
                  certificate.status === 'INACTIVE' ?
                    <Button
                      size="md"
                      variant="primary"
                      onClick={() => this.onClickCertificateAction(certificate.id, 'enable')}>
                      Enable
                    </Button> :
                    <Button
                      size="md"
                      variant="destructive"
                      onClick={() => this.onClickCertificateAction(certificate.id, 'disable')}>
                      Disable
                    </Button>
                  :''
              }
              <Button
                size="md"
                variant="destructive"
                onClick={() => this.onClickCertificateAction(certificate.id, 'delete')}>
                Delete
              </Button>
              {
                certificate.status !== 'REVOKED' &&
                <Button
                  size="md"
                  variant="destructive"
                  onClick={() => this.onClickCertificateAction(certificate.id, 'revoked')}>
                  Revoke
                </Button>
              }
            </CertificateList>
            {isCertificateModalOpen && this.createCertificateModal()}
          </>
        ))}
      </div>
    </div>
  }
}
