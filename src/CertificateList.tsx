import React, { FunctionComponent } from 'react';
import { HorizontalGroup } from '@grafana/ui';
import './ConfigEditor.scss';
import { Certificate } from './types';

interface ColumnProps {
  certificate: Certificate;
  //The content to display inside the column
  children?: React.ReactNode;
  style?: any;
}

const CertificateList: FunctionComponent<ColumnProps> = ({ certificate, children, ...rest }) => {
  return (
    <section className="card-section card-list-layout-list">
      <ol className="card-list">
        <li className="card-item-wrapper">
          <div className="card-item">
            <div className="card-item-body">
              <div className="card-item-details card-item-list-certi">
                <div className="card-item-name cared-item-certificate">
                  <div className="card-item-list-box">
                    <p>
                      <b>Topic Prefix:</b>
                      {certificate.topic}
                    </p>
                    <p>
                      <b>Client Prefix:</b>
                      {certificate.client}
                    </p>
                  </div>
                  <HorizontalGroup width="auto" className="card-item-certi-details">
                    {children}
                  </HorizontalGroup>
                </div>
              </div>
            </div>
          </div>
        </li>
      </ol>
    </section>
  );
};

export default CertificateList;
