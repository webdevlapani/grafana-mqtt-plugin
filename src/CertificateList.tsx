import React, { FunctionComponent } from "react";
import {HorizontalGroup} from "@grafana/ui";
import './ConfigEditor.scss'
import {ICertificate} from "./types";

interface IColumnProps {
  certificate: ICertificate;
  //The content to display inside the column
  children?: React.ReactNode;
  style?: any;
}

const CertificateList: FunctionComponent<IColumnProps> = ({ certificate, children, ...rest }) => {

  return <section className="card-section card-list-layout-list">
    <ol className="card-list">
      <li className="card-item-wrapper">
        <div className="card-item">
          <div className="card-item-body">
            <div className="card-item-details card-item-list-certi">
              <div className="card-item-name cared-item-certificate">
                <div className="card-item-list-box">
                  <p>Topic Prefix: {certificate.topic}</p>
                  <p>Client Prefix:{certificate.client}</p>
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
};

export default CertificateList;
